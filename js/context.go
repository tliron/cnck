package js

import (
	"context"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/logging"
	urlpkg "github.com/tliron/kutil/url"
)

//
// Context
//

type Context struct {
	Namespace  string
	Dynamic    *kubernetesutil.Dynamic
	Context    context.Context
	URLContext *urlpkg.Context
	Log        logging.Logger
}

func NewContext(namespace string, dynamic *kubernetesutil.Dynamic, context context.Context, urlContext *urlpkg.Context, log logging.Logger) *Context {
	return &Context{
		Namespace:  namespace,
		Dynamic:    dynamic,
		Context:    context,
		URLContext: urlContext,
		Log:        log,
	}
}
