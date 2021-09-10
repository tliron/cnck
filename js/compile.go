package js

import (
	"fmt"
	"strings"
)

func Compile(content string) (string, error) {
	if tags, final, err := getTags(content); err == nil {
		var builder Builder

		if len(tags) == 0 {
			// Optimize for when there are no tags
			builder.WriteLiteral(content)
		} else {
			last := 0

			for _, tag := range tags {
				// Previous chunk
				builder.WriteLiteral(content[last:tag.start])
				last = tag.end

				code := content[tag.start+2 : tag.end-2]
				trimmedCode := strings.TrimSpace(code)

				if trimmedCode == "" {
					continue
				}

				builder.Builder.WriteString(trimmedCode)
				builder.Builder.WriteRune('\n')
			}

			if last <= final {
				// Leftover chunk
				builder.WriteLiteral(content[last:])
			}
		}

		string_ := builder.Builder.String()
		//fmt.Printf("%s\n", string_)

		return string_, nil
	} else {
		return "", err
	}
}

type tag struct {
	start int
	end   int
}

func getTags(content string) ([]tag, int, error) {
	var tags []tag
	final := len(content) - 1
	start := -1

	for index, rune_ := range content {
		switch rune_ {
		case '<':
			// Opening delimiter?
			if (index < final) && (content[index+1] == '%') {
				// Not escaped?
				if (index == 0) || (content[index-1] != '\\') {
					start = index
					index += 2
				}
			}

		case '%':
			// Closing delimiter?
			if (index < final) && (content[index+1] == '>') {
				// Not escaped?
				if (index == 0) || (content[index-1] != '\\') {
					index += 2
					if start != -1 {
						tags = append(tags, tag{start, index})
						start = -1
					} else {
						return nil, -1, fmt.Errorf("closing delimiter without an opening delimiter at position %d", index)
					}
				}
			}
		}
	}

	return tags, final, nil
}
