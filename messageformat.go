// Package messageformat implements ICU messages formatting for Go.
// see http://userguide.icu-project.org/formatparse/messages
// inspired by https://github.com/SlexAxton/messageformat.js
package messageformat

import (
	"bytes"
	"fmt"
)

type MessageFormat struct {
	root      Node
	hasSelect bool
	plural    pluralFunc
}

func (x *MessageFormat) Root() Node      { return x.root }
func (x *MessageFormat) HasSelect() bool { return x.hasSelect }

func (x *MessageFormat) ToString() (string, error) {
	var buf bytes.Buffer
	err := x.root.ToString(&buf)
	return buf.String(), err
}

func (x *MessageFormat) Format() (string, error) {
	return x.FormatData(nil)
}

func (x *MessageFormat) FormatData(data Data) (string, error) {
	var buf bytes.Buffer
	err := x.root.Format(x, &buf, data, "")
	return buf.String(), err
}

func (x *MessageFormat) getNamedKey(value interface{}, ordinal bool) (string, error) {
	if nil == x.plural {
		return "", fmt.Errorf("UndefinedPluralFunc")
	}
	return x.plural(value, ordinal), nil
}
