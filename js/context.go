package js

import (
	"context"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
)

//
// Context
//

type Context struct {
	Namespace  string
	Dynamic    *kubernetesutil.Dynamic
	Context    context.Context
	URLContext *exturl.Context
	Log        commonlog.Logger
}

func NewContext(namespace string, dynamic *kubernetesutil.Dynamic, context context.Context, urlContext *exturl.Context, log commonlog.Logger) *Context {
	return &Context{
		Namespace:  namespace,
		Dynamic:    dynamic,
		Context:    context,
		URLContext: urlContext,
		Log:        log,
	}
}
