package messageformat

import (
	"bytes"
	"fmt"
	"strconv"
)

type Ordinal struct {
	Select
}

func newOrdinal(parent Node, varname string) *Ordinal {
	nd := &Ordinal{
		Select: Select{
			varname: varname,
			Choices: make(map[string]Node),
		},
	}
	parent.Add(nd)
	return nd
}

// It will returns an error if :
// - the associated value can't be convert to string or to an int (i.e. bool, ...)
// - the pluralFunc is not defined (MessageFormat.getNamedKey)
//
// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - the computed named key (MessageFormat.getNamedKey) is not a key of the given map
func (nd *Ordinal) Format(mf *MessageFormat, output *bytes.Buffer, data Data, _ string) error {
	key := nd.Varname()
	value, err := data.ValueStr(key)
	if err != nil {
		return err
	}

	var choice Node

	if v, ok := data[key]; ok {
		switch v.(type) {
		case int, float64:
		case string:
			_, err = strconv.ParseFloat(v.(string), 64)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Ordinal: Unsupported type for named key: %T", v)
		}

		key, err = mf.getNamedKey(v, true)
		if err != nil {
			return err
		}
		choice = nd.Choices[key]
	}

	if choice == nil {
		choice = nd.Other
	}
	return choice.Format(mf, output, data, value)
}
