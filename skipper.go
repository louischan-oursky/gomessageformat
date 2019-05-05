package messageformat

import (
	"fmt"

	"github.com/empirefox/makeplural/cases"
	"github.com/empirefox/makeplural/plural"
)

type SelectSkipper interface {
	Skip(key string) bool
}

type PluralSkipper struct {
	skipAll bool
	m       map[string]*cases.Case
}

func NewPluralSkipper(culture string, ordinal bool) (*PluralSkipper, error) {
	if plural.IsOthers(culture) {
		return &PluralSkipper{skipAll: true}, nil
	}

	c, ok := plural.CulturesMap()[culture]
	if !ok {
		return nil, fmt.Errorf("culture name not found from plural.Cultures: %s", culture)
	}

	var s PluralSkipper
	if ordinal {
		s.m = c.Ordinal.ToMap()
	} else {
		s.m = c.Cardinal.ToMap()
	}

	return &s, nil
}

func (s *PluralSkipper) Skip(k string) bool {
	if k == "other" {
		return false
	}

	if s.skipAll {
		return true
	}

	_, ok := s.m[k]
	return !ok
}
