package js

import (
	"context"

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

func (self *K8s) Select(arg map[string]interface{}) ([]map[string]interface{}, error) {
	var gvk schema.GroupVersionKind
	gvk.Version = "v1"

	namespace := self.namespace
	labels := make(map[string]string)

	if group, ok := arg["group"]; ok {
		gvk.Group = group.(string)
	}
	if version, ok := arg["version"]; ok {
		gvk.Version = version.(string)
	}
	if kind, ok := arg["kind"]; ok {
		gvk.Kind = kind.(string)
	}
	if namespace_, ok := arg["namespace"]; ok {
		namespace = namespace_.(string)
	}
	if labels_, ok := arg["labels"]; ok {
		labels__ := labels_.(map[string]interface{})
		for key, value := range labels__ {
			labels[key] = value.(string)
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
