package index

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

type Cfg struct {
	ID       int         `json:"_id"`
	Name     string      `json:"name"`
	CustomID string      `json:"customID,omitempty"`
	Mapping  doc.Mapping `json:"mapping"`
}

func NewCfg(tbl string, m doc.Mapping, id string) *Cfg {
	return &Cfg{
		Mapping:  m,
		CustomID: id,
		ID:       1,
		Name:     tbl,
	}
}

func DefaultCfg() *Cfg {
	return NewCfg(defaultTbl, doc.DefaultMapping(), "")
}

func NewCfgFromParams(settings string) (*Cfg, error) {
	params, err := param.Parse(settings)
	if err != nil {
		return nil, err
	}
	cfg := NewCfg(params.Index, NewMappingFromParams(params), params.UID)
	return cfg, nil
}

func NewMappingFromParams(params *param.Params) doc.Mapping {
	if !params.Has(param.SrchAttr) && !params.Has(param.FacetAttr) {
		return doc.DefaultMapping()
	}

	m := doc.NewMapping()

	for _, attr := range params.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range params.FacetAttr {
		m.AddKeywords(attr)
	}

	return m
}

func (tbl *Cfg) WithCustomID(name string) *Cfg {
	tbl.CustomID = name
	return tbl
}

func (c *Cfg) SetID(id int) {
	c.ID = id
}

func (c *Cfg) GetID() int {
	return c.ID
}

func (c *Cfg) AfterFind(db *hare.Database) error {
	//println("after find")
	return nil
}
