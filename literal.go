package messageformat

import (
	"bytes"
)

type Literal struct {
	Container
	Varname string
	Content []string
}

func (nd *Literal) Type() string { return "literal" }

func (nd *Literal) Format(_ *MessageFormat, output *bytes.Buffer, data Data, pound string) error {
	for _, c := range nd.Content {
		if c != "" {
			output.WriteString(c)
		} else if pound != "" {
			output.WriteString(pound)
		} else {
			output.WriteRune(PoundChar)
		}
	}
	return nil
}

func parseLiteral(parent Node, start, end int, input []rune) {
	var items []int

	escaped := false

	s, e := start, start
	gap := 0
	for i := start; i < end; i++ {
		c := input[i]

		if EscapeChar == c {
			gap++
			e++
			escaped = true
		} else {
			switch c {
			default:
				e++

			case OpenChar, CloseChar, PoundChar:
				if escaped {
					if i-s > gap {
						if gap > 1 {
							items = append(items, s, i)
						} else {
							items = append(items, s, i-1)
						}
					}
					s = i
				} else {
					if s != e {
						items = append(items, s, e, i, i)
					} else if s != i {
						items = append(items, s, i, i, i)
					} else {
						items = append(items, i, i)
					}
					s = i + 1
				}
				e = s
			}

			escaped = false
			gap = 0
		}
	}

	if s < end {
		items = append(items, s, end)
	}

	n := len(items)
	content := make([]string, n/2)
	for i := 0; i < n; i += 2 {
		content[i/2] = string(input[items[i]:items[i+1]])
	}

	child := &Literal{Content: content}
	if sn, ok := parent.(Varnamer); ok {
		child.Varname = sn.Varname()
	}
	parent.Add(child)
}
