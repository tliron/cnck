package controller

import (
	contextpkg "context"
	"time"

	"github.com/tliron/commonlog"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	core "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	dynamicpkg "k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/record"
)

//
// Controller
//

type Controller struct {
	Namespace   string
	Kubernetes  kubernetes.Interface
	Dynamic     *kubernetesutil.Dynamic
	StopChannel <-chan struct{}

	Processors *kubernetesutil.Processors
	Events     record.EventRecorder

	KubernetesInformerFactory informers.SharedInformerFactory

	ConfigMaps listers.ConfigMapLister

	Context contextpkg.Context
	Log     commonlog.Logger
}

func NewController(toolName string, cluster bool, namespace string, kubernetes kubernetes.Interface, dynamic dynamicpkg.Interface, informerResyncPeriod time.Duration, stopChannel <-chan struct{}) *Controller {
	context := contextpkg.TODO()

	if cluster {
		namespace = ""
	}

	log := commonlog.GetLoggerf("%s.controller", toolName)

	self := Controller{
		Namespace:   namespace,
		Kubernetes:  kubernetes,
		Dynamic:     kubernetesutil.NewDynamic(toolName, dynamic, kubernetes.Discovery(), namespace, context),
		StopChannel: stopChannel,
		Processors:  kubernetesutil.NewProcessors(toolName),
		Events:      kubernetesutil.CreateEventRecorder(kubernetes, "CNCK", log),
		Context:     context,
		Log:         log,
	}

	if cluster {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactory(kubernetes, informerResyncPeriod)
	} else {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactoryWithOptions(kubernetes, informerResyncPeriod, informers.WithNamespace(namespace))
	}

	// Informers
	configMapInformer := self.KubernetesInformerFactory.Core().V1().ConfigMaps()

	// Listers
	self.ConfigMaps = configMapInformer.Lister()

	// Processors

	processorPeriod := 5 * time.Second

	self.Processors.Add(ConfigMapGVK, kubernetesutil.NewProcessor(
		toolName,
		"configmaps",
		configMapInformer.Informer(),
		processorPeriod,
		func(name string, namespace string) (interface{}, error) {
			return self.GetConfigMap(namespace, name)
		},
		func(object interface{}) (bool, error) {
			return self.processConfigMap(object.(*core.ConfigMap))
		},
	))

	return &self
}

func (self *Controller) Run(concurrency uint, startup func()) error {
	defer utilruntime.HandleCrash()

	self.Log.Info("starting informer factories")
	self.KubernetesInformerFactory.Start(self.StopChannel)

	self.Log.Info("waiting for processor informer caches to sync")
	utilruntime.HandleError(self.Processors.WaitForCacheSync(self.StopChannel))

	self.Log.Infof("starting processors (concurrency=%d)", concurrency)
	self.Processors.Start(concurrency, self.StopChannel)
	defer self.Processors.ShutDown()

	if startup != nil {
		go startup()
	}

	<-self.StopChannel

	self.Log.Info("shutting down")

	return nil
}
