package messageformat

import (
	"bytes"
	"errors"
	"sort"
)

type Varnamer interface {
	Varname() string
}

type selectParser interface {
	sort.Interface
	parse(p *Parser, skipper SelectSkipper,
		char rune, start, end int, input []rune) (int, error)
}

type Choice struct {
	Key  string
	Node Node
}

type Select struct {
	Container
	varname    string
	ChoicesMap map[string]Node
	Choices    []Choice
	Other      Node
}

func newSelect(parent Node, varname string) *Select {
	nd := &Select{
		varname:    varname,
		ChoicesMap: make(map[string]Node, 5),
		Choices:    make([]Choice, 0, 5),
	}
	parent.Add(nd)
	return nd
}

// sort Choices
func (s *Select) Len() int {
	return len(s.Choices)
}

// sort Choices
func (s *Select) Swap(i, j int) {
	s.Choices[i], s.Choices[j] = s.Choices[j], s.Choices[i]
}

// sort Choices
func (s *Select) Less(i, j int) bool {
	return s.Choices[i].Key < s.Choices[j].Key
}

func (nd *Select) addChoice(key string, choice Node) {
	nd.ChoicesMap[key] = choice
	nd.Choices = append(nd.Choices, Choice{Key: key, Node: choice})
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
			nd.Other = choice
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

	if nd.Other == nil {
		return pos, errors.New("MissingMandatoryChoice")
	}
	return pos, nil
}

func (nd *Select) Varname() string { return nd.varname }
func (nd *Select) Type() string    { return "select" }

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
		choice = nd.Other
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
