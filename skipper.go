package messageformat

import (
	"fmt"

	"github.com/empirefox/makeplural/plural"
)

type SelectSkipper interface {
	Skip(key string) bool
}

type PluralSkipper struct {
	skipNotPluralForms bool
	m                  map[string]string
}

func NewPluralSkipper(culture string, ordinal, skipNotPluralForms bool) (*PluralSkipper, error) {
	c, ok := plural.Cultures[culture]
	if !ok {
		return nil, fmt.Errorf("culture name not found from plural.Cultures: %s", culture)
	}

	s := PluralSkipper{skipNotPluralForms: skipNotPluralForms}
	if ordinal {
		s.m = c.Ordinal
	} else {
		s.m = c.Cardinal
	}

	return &s, nil
}

func (s *PluralSkipper) Skip(k string) bool {
	if k == "other" {
		return false
	}
	_, ok := s.m[k]
	if s.skipNotPluralForms {
		return !ok
	}
	return ok
}
