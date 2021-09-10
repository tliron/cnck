package controller

import (
	"strings"

	"github.com/tliron/cnck/js"
	"github.com/tliron/kutil/logging"
	urlpkg "github.com/tliron/kutil/url"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const annotationRender = "cnck.github.com/render"
const annotationRefresh = "cnck.github.com/refresh"
const templateSuffix = ".template"

func (self *Controller) GetConfigMap(namespace string, networkName string) (*core.ConfigMap, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	if configMap, err := self.Kubernetes.CoreV1().ConfigMaps(namespace).Get(self.Context, networkName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if configMap.Kind == "" {
			configMap = configMap.DeepCopy()
			configMap.APIVersion, configMap.Kind = ConfigMapGVK.ToAPIVersionAndKind()
		}
		return configMap, nil
	} else {
		return nil, err
	}
}

func (self *Controller) processConfigMap(configMap *core.ConfigMap) (bool, error) {
	if templateName, ok := configMap.Annotations[annotationRender]; ok {
		if strings.HasSuffix(templateName, templateSuffix) {
			if template, ok := configMap.Data[templateName]; ok {
				if script, err := js.Compile(template); err == nil {
					renderedName := templateName[:len(templateName)-len(templateSuffix)]

					urlContext := urlpkg.NewContext()
					defer urlContext.Release()
					context := js.NewContext(self.Namespace, self.Dynamic, self.Context, urlContext, logging.NewSubLogger(self.Log, "js"))

					if rendered, err := context.Render(script); err == nil {
						self.Log.Infof("rendered: %s", renderedName)

						exists := false
						if existingRendered, ok := configMap.Data[renderedName]; ok {
							if existingRendered == rendered {
								exists = true
							}
						}

						if exists {
							self.Log.Info("no change to ConfigMap")
						} else {
							configMap_ := configMap.DeepCopy()
							configMap_.Data[renderedName] = rendered

							if _, err := self.Kubernetes.CoreV1().ConfigMaps(configMap.Namespace).Update(self.Context, configMap_, meta.UpdateOptions{}); err == nil {
								self.Log.Info("ConfigMap updated")
							} else {
								return false, err
							}
						}
					} else {
						self.Log.Infof("could not render: %s", renderedName)
						return false, err
					}
				} else {
					self.Log.Infof("could not compile template: %s", templateName)
					return true, err
				}
			} else {
				self.Log.Warningf("template does not exist: %s", templateName)
			}
		} else {
			self.Log.Warningf("%q annotation does not have %q suffix: %s", annotationRender, templateSuffix, templateName)
		}
	}

	return true, nil
}
