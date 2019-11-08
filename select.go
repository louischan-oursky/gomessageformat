package messageformat

import (
	"bytes"
	"errors"
	"sort"
	"strconv"
)

type SelectNode interface {
	Node
	sort.Interface
	Varname() string
	Choices() []Choice
	ChoicesUnsorted() []Choice
	Other() Node
}

type selectParser interface {
	SelectNode
	parse(p *Parser, skipper SelectSkipper,
		char rune, start, end int, input []rune) (int, error)
}

type Choice struct {
	Key  string
	Node Node
}

type Select struct {
	Container
	varname         string
	ChoicesMap      map[string]Node
	choicesUnsorted []Choice
	choices         []Choice
	other           Node
}

func newSelect(parent Node, varname string) *Select {
	nd := &Select{
		varname:         varname,
		ChoicesMap:      make(map[string]Node, 5),
		choices:         make([]Choice, 0, 5),
		choicesUnsorted: make([]Choice, 0, 5),
	}
	AddChild(parent, nd)
	return nd
}

// sort choices
func (s *Select) Len() int {
	return len(s.choices)
}

// sort choices
func (s *Select) Swap(i, j int) {
	s.choices[i], s.choices[j] = s.choices[j], s.choices[i]
}

// sort choices
func (s *Select) Less(i, j int) bool {
	return s.choices[i].Key < s.choices[j].Key
}

func (nd *Select) addChoice(key string, choice Node) {
	nd.ChoicesMap[key] = choice
	nd.choices = append(nd.choices, Choice{Key: key, Node: choice})
}

func (nd *Select) parse(p *Parser, skipper SelectSkipper,
	char rune, start, end int, input []rune) (int, error) {
	if char != PartChar {
		return start, errors.New("MalformedOption")
	}

	pos := start + 1

	for pos < end {
		key, char, i, err := readKey(char, pos, end, input)

		if err != nil {
			return i, err
		}
		if char == ':' {
			return i, errors.New("UnexpectedExtension")
		}

		choice, char, i, err := p.readChoice(nd, char, i, end, input)
		if err != nil {
			return i, err
		}

		if key == "other" {
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

		pos = i

		if char == CloseChar {
			break
		}
	}

	if nd.other == nil {
		return pos, errors.New("MissingMandatoryChoice")
	}

	nd.choicesUnsorted = append(nd.choicesUnsorted, nd.choices...)
	return pos, nil
}

func (nd *Select) Varname() string           { return nd.varname }
func (nd *Select) Type() string              { return "select" }
func (nd *Select) Choices() []Choice         { return nd.choices }
func (nd *Select) ChoicesUnsorted() []Choice { return nd.choicesUnsorted }
func (nd *Select) Other() Node               { return nd.other }

func (nd *Select) ToString(output *bytes.Buffer) error {
	return selectNodeToString(nd, output)
}

func selectNodeToString(nd SelectNode, output *bytes.Buffer) (err error) {
	_, err = output.WriteRune(OpenChar)
	if err != nil {
		return
	}

	_, err = output.WriteString(nd.Varname())
	if err != nil {
		return
	}

	_, err = output.WriteString(`, `)
	if err != nil {
		return
	}

	_, err = output.WriteString(nd.Type())
	if err != nil {
		return
	}

	_, err = output.WriteString(`, `)
	if err != nil {
		return
	}

	if s, ok := nd.(equals); ok {
		if offset := s.Offset(); offset != 0 {
			_, err = output.WriteString(`offset:`)
			if err != nil {
				return
			}
			_, err = output.WriteString(strconv.Itoa(offset))
			if err != nil {
				return
			}
			_, err = output.WriteRune(' ')
			if err != nil {
				return
			}
		}
		for _, eq := range s.Equals() {
			_, err = output.WriteRune(EqualChar)
			if err != nil {
				return
			}
			_, err = output.WriteString(strconv.FormatFloat(eq.Key, 'f', -1, 64))
			if err != nil {
				return
			}
			_, err = output.WriteRune(OpenChar)
			if err != nil {
				return
			}
			err = eq.Node.ToString(output)
			if err != nil {
				return
			}
			_, err = output.WriteRune(CloseChar)
			if err != nil {
				return
			}
		}
	}

	for _, choice := range nd.ChoicesUnsorted() {
		_, err = output.WriteString(choice.Key)
		if err != nil {
			return
		}
		_, err = output.WriteRune(OpenChar)
		if err != nil {
			return
		}
		err = choice.Node.ToString(output)
		if err != nil {
			return
		}
		_, err = output.WriteRune(CloseChar)
		if err != nil {
			return
		}
	}

	_, err = output.WriteString(`other`)
	if err != nil {
		return
	}
	_, err = output.WriteRune(OpenChar)
	if err != nil {
		return
	}
	err = nd.Other().ToString(output)
	if err != nil {
		return
	}
	_, err = output.WriteRune(CloseChar)
	if err != nil {
		return
	}

	_, err = output.WriteRune(CloseChar)
	return
}

// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - its string representation is not a key of the given map
//
// It will returns an error if :
// - the associated value can't be convert to string (i.e. bool, ...)
func (nd *Select) Format(mf *MessageFormat, output *bytes.Buffer, data Data, _ string) error {
	value, err := data.ValueStr(nd.Varname())
	if err != nil {
		return err
	}

	choice, ok := nd.ChoicesMap[value]
	if !ok {
		choice = nd.other
	}
	return choice.Format(mf, output, data, value)
}

func readKey(char rune, start, end int, input []rune) (string, rune, int, error) {
	char, pos := whitespace(start, end, input)
	fc_pos, lc_pos := pos, pos

	for pos < end {
		switch char {
		default:
			lc_pos = pos + 1

		case ' ', '\r', '\n', '\t':
			char, pos = whitespace(pos+1, end, input)
			return string(input[fc_pos:lc_pos]), char, pos, nil

		case ':', PartChar, CloseChar, OpenChar:
			if fc_pos != lc_pos {
				return string(input[fc_pos:lc_pos]), char, pos, nil
			}
			return "", char, pos, errors.New("MissingChoiceName")
		}

		pos++

		if pos < end {
			char = input[pos]
		}
	}
	return "", char, pos, errors.New("UnbalancedBraces")
}

func (p *Parser) readChoice(parent Node, char rune, pos, end int, input []rune) (*Container, rune, int, error) {
	if char != OpenChar {
		return nil, char, pos, errors.New("MissingChoiceContent")
	}

	choice := newContainer(parent)
	pos, _, err := p.parse(pos+1, end, input, choice)
	if err != nil {
		return nil, char, pos, err
	}

	pos++
	if pos < end {
		char = input[pos]
	}

	if isWhitespace(char) {
		char, pos = whitespace(pos+1, end, input)
	}
	return choice, char, pos, nil
}
