package js

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

func (self *Context) NewEnvironment(builder *strings.Builder, scriptlet string) *commonjs.Environment {
	environment := commonjs.NewEnvironment(self.URLContext, nil)

	environment.CreateResolver = func(url exturl.URL, context *commonjs.Context) commonjs.ResolveFunc {
		return func(id string, raw bool) (exturl.URL, error) {
			url := self.URLContext.NewInternalURL(id)
			url.Content = util.StringToBytes(scriptlet)
			return url, nil
		}
	}

	environment.Extensions = []commonjs.Extension{
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

func (self *Context) createLogExtension(context *commonjs.Context) goja.Value {
	return context.Environment.Runtime.ToValue(self.Log)
}

func (self *Context) createK8sExtension(context *commonjs.Context) goja.Value {
	k8s := K8s{
		namespace: self.Namespace,
		dynamic:   self.Dynamic,
		context:   self.Context,
	}
	return context.Environment.Runtime.ToValue(&k8s)
}

func createWriteExtension(builder *strings.Builder) commonjs.CreateExtensionFunc {
	write := func(arg interface{}) error {
		_, err := builder.WriteString(fmt.Sprintf("%s", arg))
		return err
	}

	return func(context *commonjs.Context) goja.Value {
		return context.Environment.Runtime.ToValue(write)
	}
}
