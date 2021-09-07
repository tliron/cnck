package js

import (
	"strings"
)

//
// Builder
//

type Builder struct {
	strings.Builder
}

func (self *Builder) WriteLiteral(literal string) {
	if literal != "" {
		self.WriteString("write('")
		for _, rune_ := range literal {
			switch rune_ {
			case '\n':
				self.WriteString("\\n")
			case '\'':
				self.WriteString("\\'")
			case '\\':
				self.WriteString("\\\\")
			default:
				self.WriteRune(rune_)
			}
		}
		self.WriteString("');\n")
	}
}
