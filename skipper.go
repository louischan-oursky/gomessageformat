package messageformat

import (
	"fmt"

	"github.com/empirefox/makeplural/plural"
	"golang.org/x/text/language"
)

var pluralForms = map[string]int{
	"zero":  0,
	"one":   1,
	"two":   2,
	"few":   3,
	"many":  4,
	"other": 5,
}

type SelectSkipper interface {
	Skip(key string) (skip bool, err error)
}

type NothingSkipper struct{}

func (s NothingSkipper) Skip(k string) (skip bool, err error) {
	return false, nil
}

type PluralSkipper struct {
	skipAll bool
	m       map[string]*plural.Case
}

func NewPluralSkipper(culture language.Tag, ordinal bool) (*PluralSkipper, error) {
	if plural.Info.IsOthers(culture) {
		return &PluralSkipper{skipAll: true}, nil
	}

	c, ok := plural.Info.CulturesMap()[culture]
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

func (s *PluralSkipper) Skip(k string) (skip bool, err error) {
	_, ok := pluralForms[k]
	if !ok {
		return false, fmt.Errorf("plural form not found: %s", k)
	}

	if k == "other" {
		return false, nil
	}

	if s.skipAll {
		return true, nil
	}

	_, ok = s.m[k]
	return !ok, nil
}
