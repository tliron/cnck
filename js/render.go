package js

import (
	"strings"
)

func (self *Context) Render(script string) (string, error) {
	var builder strings.Builder
	environment := self.NewEnvironment(&builder, script)
	defer environment.Release()

	_, err := environment.RequireID("cnck.script")
	return builder.String(), err
}
