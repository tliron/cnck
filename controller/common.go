package controller

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var ConfigMapGVK = schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}
