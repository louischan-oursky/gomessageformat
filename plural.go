package messageformat

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type Equal struct {
	Key  float64
	Node Node
}

type Equals []Equal

func (s Equals) Len() int           { return len(s) }
func (s Equals) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Equals) Less(i, j int) bool { return s[i].Key < s[j].Key }

type equals interface {
	Equals() Equals
	Offset() int
}

type Plural struct {
	Select
	EqualsMap map[float64]Node
	equals    Equals
	offset    int
}

func newPlural(parent Node, varname string) *Plural {
	nd := &Plural{
		Select: Select{
			varname:    varname,
			ChoicesMap: make(map[string]Node, 5),
			choices:    make([]Choice, 0, 5),
		},
		EqualsMap: make(map[float64]Node, 0),
		equals:    make(Equals, 0, 5),
	}
	AddChild(parent, nd)
	return nd
}

func (nd *Plural) Type() string   { return "plural" }
func (nd *Plural) Equals() Equals { return nd.equals }
func (nd *Plural) Offset() int    { return nd.offset }

// sort choices
func (s *Plural) Less(i, j int) bool {
	return pluralForms[s.choices[i].Key] < pluralForms[s.choices[j].Key]
}

func (nd *Plural) addEqual(f64 float64, choice Node) {
	nd.EqualsMap[f64] = choice
	nd.equals = append(nd.equals, Equal{Key: f64, Node: choice})
}

func (nd *Plural) parse(p *Parser, skipper SelectSkipper,
	char rune, start, end int, input []rune) (int, error) {
	if char != PartChar {
		return start, fmt.Errorf("MalformedOption")
	}

	pos := start + 1

	for pos < end {
		key, char, i, err := readKey(char, pos, end, input)

		if err != nil {
			return i, err
		}

		if char == ':' {
			if key != "offset" {
				return i, fmt.Errorf("UnsupportedExtension: `%s`", key)
			}

			offset, c, j, err := readOffset(i+1, end, input)
			if err != nil {
				return j, err
			}

			nd.offset = offset

			if isWhitespace(c) {
				j++
			}

			k, c, j, err := readKey(c, j, end, input)

			if k == "" {
				return j, fmt.Errorf("MissingChoiceName")
			}

			key, char, i = k, c, j
		}

		choice, c, i, err := p.readChoice(nd, char, i, end, input)
		if err != nil {
			return i, err
		}

		if key[0] == EqualChar {
			f64, err := strconv.ParseFloat(key[1:], 64)
			if err != nil {
				return i, fmt.Errorf("invalid number key `%s`", key)
			}
			nd.addEqual(f64, choice)
		} else if key == "other" {
			nd.other = choice
		} else {
			skip, err := skipper.Skip(key)
			if err != nil {
				return i, err
			}
			if !skip {
				nd.addChoice(key, choice)
			}
		}
		pos, char = i, c

		if CloseChar == char {
			break
		}
	}

	if nd.other == nil {
		return pos, errors.New("MissingMandatoryChoice")
	}

	nd.choicesUnsorted = append(nd.choicesUnsorted, nd.choices...)
	sort.Sort(nd.equals)
	return pos, nil
}

func (nd *Plural) ToString(output *bytes.Buffer) error {
	return selectNodeToString(nd, output)
}

// It will returns an error if :
// - the associated value can't be convert to string or to an int (i.e. bool, ...)
// - the pluralFunc is not defined (MessageFormat.getNamedKey)
//
// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - the computed named key (MessageFormat.getNamedKey) is not a key of the given map
func (nd *Plural) Format(mf *MessageFormat, output *bytes.Buffer, data Data, _ string) error {
	key := nd.Varname()
	offset := nd.offset

	value, err := data.ValueStr(key)
	if err != nil {
		return err
	}

	var choice Node

	if iv, ok := data[key]; ok {
		switch v := iv.(type) {
		case int:
			choice = nd.EqualsMap[float64(v)]
		case float64:
			choice = nd.EqualsMap[v]
		case string:
			f64, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("Plural: value not a number string: %s", iv)
			}
			choice = nd.EqualsMap[f64]
		default:
			return fmt.Errorf("Plural: Unsupported type for named key: %T", iv)
		}

		if choice == nil {
			switch iv.(type) {
			case int:
				if 0 != offset {
					offset_value := iv.(int) - offset
					value = fmt.Sprintf("%d", offset_value)
					key, err = mf.getNamedKey(offset_value, false)
				} else {
					key, err = mf.getNamedKey(iv.(int), false)
				}

			case float64:
				if 0 != offset {
					offset_value := iv.(float64) - float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = mf.getNamedKey(offset_value, false)
				} else {
					key, err = mf.getNamedKey(iv.(float64), false)
				}

			case string:
				if 0 != offset {
					offset_value, fError := strconv.ParseFloat(value, 64)
					if nil != fError {
						return fError
					}
					offset_value -= float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = mf.getNamedKey(offset_value, false)
				} else {
					key, err = mf.getNamedKey(value, false)
				}
			}

			if err != nil {
				return err
			}
			choice = nd.ChoicesMap[key]
		}
	}

	if choice == nil {
		choice = nd.other
	}
	return choice.Format(mf, output, data, value)
}

func readOffset(start, end int, input []rune) (int, rune, int, error) {
	var buf bytes.Buffer
	char, pos := whitespace(start, end, input)

	for pos < end {
		switch char {
		default:
			buf.WriteRune(char)
			pos++

			if pos < end {
				char = input[pos]
			}

		case ' ', '\r', '\n', '\t', OpenChar, CloseChar:
			if 0 != buf.Len() {
				result, err := strconv.Atoi(buf.String())
				if err != nil {
					return 0, char, pos, fmt.Errorf("BadCast")
				}
				if result < 0 {
					return 0, char, pos, fmt.Errorf("InvalidOffsetValue")
				}
				return result, char, pos, nil
			}
			return 0, char, pos, fmt.Errorf("MissingOffsetValue")
		}
	}
	return 0, char, pos, fmt.Errorf("UnbalancedBraces")
}
