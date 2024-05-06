package doc

import (
	"github.com/ohzqq/srch/analyzer"
)

type Mapping map[analyzer.Analyzer][]string

func NewMapping() Mapping {
	return make(map[analyzer.Analyzer][]string)
}

func DefaultMapping() Mapping {
	m := NewMapping()
	m.AddSimple("*")
	return m
}

func (m Mapping) AddFulltext(name ...string) {
	m[analyzer.Standard] = append(m[analyzer.Standard], name...)
}

func (m Mapping) AddKeywords(name ...string) {
	m[analyzer.Keyword] = append(m[analyzer.Keyword], name...)
}

func (m Mapping) AddSimple(name ...string) {
	m[analyzer.Simple] = append(m[analyzer.Simple], name...)
}
