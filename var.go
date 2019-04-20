package messageformat

import (
	"bytes"
	"errors"
)

type Var struct {
	Container
	Varname string
}

func (nd *Var) Type() string { return "var" }

func (nd *Var) Format(_ *MessageFormat, output *bytes.Buffer, data Data, value string) (err error) {
	if value == "" {
		value, err = data.ValueStr(nd.Varname)
		if err != nil {
			return err
		}
	}
	_, err = output.WriteString(value)
	return err
}

func readVar(start, end int, input []rune) (string, rune, int, error) {
	char, pos := whitespace(start, end, input)
	fc_pos, lc_pos := pos, pos

	for pos < end {
		switch char {
		default:
			// [_0-9a-zA-Z]+
			if '_' != char && (char < '0' || char > '9') && (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
				return "", char, pos, errors.New("InvalidFormat")
			} else if pos != lc_pos { // non continu (inner whitespace)
				return "", char, pos, errors.New("InvalidFormat")
			}

			lc_pos = pos + 1

			pos++

			if pos < end {
				char = input[pos]
			}

		case ' ', '\r', '\n', '\t':
			char, pos = whitespace(pos+1, end, input)

		case PartChar, CloseChar:
			return string(input[fc_pos:lc_pos]), char, pos, nil

		case OpenChar:
			return "", char, pos, errors.New("InvalidExpr")
		}
	}
	return "", char, pos, errors.New("UnbalancedBraces")
}
