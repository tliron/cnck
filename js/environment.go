package js

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	urlpkg "github.com/tliron/exturl"
	"github.com/tliron/kutil/js"
	"github.com/tliron/kutil/util"
)

func (self *Context) NewEnvironment(builder *strings.Builder, scriptlet string) *js.Environment {
	environment := js.NewEnvironment(self.URLContext, nil)

	environment.CreateResolver = func(url urlpkg.URL, context *js.Context) js.ResolveFunc {
		return func(id string, raw bool) (urlpkg.URL, error) {
			url := urlpkg.NewInternalURL(id, self.URLContext)
			url.Content = util.StringToBytes(scriptlet)
			return url, nil
		}
	}

	environment.Extensions = []js.Extension{
		{
			Name:   "k8s",
			Create: self.createK8sExtension,
		},
		{
			Name:   "log",
			Create: self.createLogExtension,
		},
		{
			Name:   "write",
			Create: createWriteExtension(builder),
		},
	}

	return environment
}

func (self *Context) createLogExtension(context *js.Context) goja.Value {
	return context.Environment.Runtime.ToValue(self.Log)
}

func (self *Context) createK8sExtension(context *js.Context) goja.Value {
	k8s := K8s{
		namespace: self.Namespace,
		dynamic:   self.Dynamic,
		context:   self.Context,
	}
	return context.Environment.Runtime.ToValue(&k8s)
}

func createWriteExtension(builder *strings.Builder) js.CreateExtensionFunc {
	write := func(arg interface{}) error {
		_, err := builder.WriteString(fmt.Sprintf("%s", arg))
		return err
	}

	return func(context *js.Context) goja.Value {
		return context.Environment.Runtime.ToValue(write)
	}
}
