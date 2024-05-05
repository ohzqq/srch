package doc

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/param"
)

type Mapping struct {
	ID      int                            `json:"id"`
	Mapping map[analyzer.Analyzer][]string `json:"mapping"`
}

func NewMapping() *Mapping {
	return &Mapping{
		ID:      0,
		Mapping: make(map[analyzer.Analyzer][]string),
	}
}

func DefaultMapping() *Mapping {
	m := NewMapping()
	m.AddSimple("*")
	return m
}

func NewMappingFromParams(params *param.Params) *Mapping {
	m := NewMapping()

	for _, attr := range params.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range params.Facets {
		m.AddKeywords(attr)
	}

	return m
}

func (m *Mapping) AddFulltext(name ...string) {
	m.Mapping[analyzer.Standard] = append(m.Mapping[analyzer.Standard], name...)
}

func (m *Mapping) AddKeywords(name ...string) {
	m.Mapping[analyzer.Keyword] = append(m.Mapping[analyzer.Keyword], name...)
}

func (m *Mapping) AddSimple(name ...string) {
	m.Mapping[analyzer.Simple] = append(m.Mapping[analyzer.Simple], name...)
}

func (m *Mapping) SetID(id int) {
	m.ID = id
}

func (m *Mapping) GetID() int {
	return m.ID
}

func (m *Mapping) AfterFind(_ *hare.Database) error {
	return nil
}
