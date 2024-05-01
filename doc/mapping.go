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
	m := make(Mapping)

	for _, attr := range params.SrchAttr {
		m[analyzer.Fulltext] = append(m[analyzer.Fulltext], attr)
	}

	for _, attr := range params.Facets {
		m[analyzer.Keywords] = append(m[analyzer.Keywords], attr)
	}

	return m
}

func (m Mapping) AddFulltext(name ...string) {
	m[analyzer.Fulltext] = append(m[analyzer.Fulltext], name...)
}

func (m Mapping) AddKeywords(name ...string) {
	m[analyzer.Keywords] = append(m[analyzer.Keywords], name...)
}

func (m Mapping) AddSimple(name ...string) {
	m[analyzer.Simple] = append(m[analyzer.Simple], name...)
}
