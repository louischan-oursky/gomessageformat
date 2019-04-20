// Package messageformat implements ICU messages formatting for Go.
// see http://userguide.icu-project.org/formatparse/messages
// inspired by https://github.com/SlexAxton/messageformat.js
package messageformat

import (
	"bytes"
	"fmt"
)

type MessageFormat struct {
	root   Node
	plural pluralFunc
}

func (x *MessageFormat) Root() Node { return x.root }

func (x *MessageFormat) Format() (string, error) {
	return x.FormatData(nil)
}

func (x *MessageFormat) FormatData(data Data) (string, error) {
	var buf bytes.Buffer
	err := x.root.Format(x, &buf, data, "")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (x *MessageFormat) getNamedKey(value interface{}, ordinal bool) (string, error) {
	if nil == x.plural {
		return "", fmt.Errorf("UndefinedPluralFunc")
	}
	return x.plural(value, ordinal), nil
}
