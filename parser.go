package messageformat

import (
	"fmt"
	"sort"

	"github.com/louischan-oursky/gomakeplural/plural"
	"golang.org/x/text/language"
)

const (
	EscapeChar = '\\'
	OpenChar   = '{'
	CloseChar  = '}'
	PartChar   = ','
	PoundChar  = '#'
	EqualChar  = '='
)

// pluralFunc describes a function used to produce a named key
// when processing a plural or selectordinal node.
type pluralFunc func(interface{}, bool) string

type Parser struct {
	culture language.Tag
	plural  pluralFunc
	noSkip  bool
}

// SkipForPlural must be called before Parse.
func (x *Parser) SkipForPlural(skip bool) {
	x.noSkip = !skip
}

func (x *Parser) Parse(input string, data interface{}) (*MessageFormat, error) {
	runes := []rune(input)
	pos, end := 0, len(runes)

	var hasSelect bool

	root := &Root{culture: x.culture, data: data}
	for pos < end {
		i, level, err := x.parse(pos, end, runes, root, &hasSelect)
		if err != nil {
			return nil, parseError{err.Error(), i}
		} else if 0 != level {
			return nil, parseError{"UnbalancedBraces", i}
		}

		pos = i
	}
	return &MessageFormat{root, hasSelect, x.plural}, nil
}

func (x *Parser) parseNode(parent Node, start, end int, input []rune, hasSelect *bool) (int, error) {
	varname, char, pos, err := readVar(start, end, input)
	if err != nil {
		return pos, err
	}
	if varname == "" {
		return pos, fmt.Errorf("MissingVarName")
	}
	if char == CloseChar {
		AddChild(parent, &Var{Varname: varname})
		return pos, nil
	}

	ctype, char, pos, err := readVar(pos+1, end, input)
	if err != nil {
		return pos, err
	}

	var sp selectParser
	var skipper SelectSkipper = NothingSkipper{}
	switch ctype {
	case "select":
		sp = newSelect(parent, varname)
	case "selectordinal":
		sp = newOrdinal(parent, varname)
		if !x.noSkip {
			skipper, err = NewPluralSkipper(x.culture, true)
		}
	case "plural":
		sp = newPlural(parent, varname)
		if !x.noSkip {
			skipper, err = NewPluralSkipper(x.culture, false)
		}
	default:
		return pos, fmt.Errorf("UnknownType: `%s`", ctype)
	}
	if err != nil {
		return pos, err
	}

	pos, err = sp.parse(x, skipper, char, pos, end, input)
	if err != nil {
		return pos, err
	}

	if pos >= end || input[pos] != CloseChar {
		return pos, fmt.Errorf("UnbalancedBraces")
	}

	if hasSelect != nil {
		*hasSelect = true
	}
	sort.Sort(sp)
	return pos, nil
}

func (x *Parser) parse(start, end int, input []rune, nd Node, hasSelect *bool) (int, int, error) {
	pos := start
	level := 0
	escaped := false

loop:
	for pos < end {
		char := input[pos]

		switch char {
		default:
			pos++
			escaped = false

		case EscapeChar:
			pos++
			escaped = true

		case CloseChar:
			if !escaped {
				level--
				break loop
			}
			pos++
			escaped = false

		case OpenChar:
			if !escaped {
				level++

				if pos > start {
					parseLiteral(nd, start, pos, input)
				}

				i, err := x.parseNode(nd, pos+1, end, input, hasSelect)
				if err != nil {
					return i, level, err
				}

				level--

				pos = i
				start = pos + 1
			}

			pos++
			escaped = false
		}
	}

	if pos > start {
		parseLiteral(nd, start, pos, input)
	}
	return pos, level, nil
}

func NewWithCulture(culture language.Tag) (*Parser, error) {
	fn, err := plural.GetFunc(culture)
	if err != nil {
		return nil, err
	}
	return &Parser{culture: culture, plural: fn}, nil
}

func New() (*Parser, error) {
	return NewWithCulture(language.English)
}
