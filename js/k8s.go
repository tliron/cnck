package js

import (
	"context"
	"fmt"

	"github.com/tliron/go-ard"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//
// K8s
//

type K8s struct {
	namespace string
	dynamic   *kubernetesutil.Dynamic
	context   context.Context
}

func (self *K8s) Select(config ard.StringMap) ([]ard.StringMap, error) {
	var gvk schema.GroupVersionKind
	var namespace string
	labels := make(map[string]string)

	node := ard.NewNode(config)
	gvk.Group, _ = node.Get("group").String()
	gvk.Version, _ = node.Get("version").String()
	if gvk.Version == "" {
		gvk.Version = "v1"
	}
	gvk.Kind, _ = node.Get("kind").String()
	namespace, _ = node.Get("namespace").String()
	if namespace == "" {
		namespace = self.namespace
	}

	if labels_, ok := node.Get("labels").StringMap(); ok {
		for key, value := range labels_ {
			labels[key] = fmt.Sprintf("%s", value)
		}
	}

	if unstructureds, err := self.dynamic.ListResources(gvk, namespace, labels); err == nil {
		r := make([]map[string]interface{}, len(unstructureds))
		for index, unstructured := range unstructureds {
			r[index] = unstructured.Object
		}
		//format.PrintYAML(r, os.Stdout, false, false)
		return r, nil
	} else {
		return nil, err
	}
}
