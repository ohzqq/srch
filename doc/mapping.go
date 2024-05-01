package doc

import (
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/param"
)

type Mapping map[analyzer.Analyzer][]string

func NewMapping() Mapping {
	return make(Mapping)
}

func NewMappingFromParams(params *param.Params) Mapping {
	m := NewMapping()

	for _, attr := range params.SrchAttr {
		m[analyzer.Standard] = append(m[analyzer.Standard], attr)
	}

	for _, attr := range params.Facets {
		m[analyzer.Keyword] = append(m[analyzer.Keyword], attr)
	}

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
