package index

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

type IdxCfg struct {
	ID      int         `json:"_id"`
	Mapping doc.Mapping `json:"mapping"`

	*param.Cfg
}

func NewCfg() *IdxCfg {
	cfg := &IdxCfg{
		Cfg: param.NewCfg(),
	}

	return cfg.
		SetMapping(doc.DefaultMapping()).
		SetName(defaultTbl)
}

func (cfg *IdxCfg) Parse(v any) error {
	err := param.Decode(v, cfg.Cfg)
	if err != nil {
		return err
	}

	cfg.SetName(cfg.Index).
		SetMapping(NewMappingFromParamCfg(cfg.Cfg)).
		SetCustomID(cfg.UID)

	return nil
}

func (cfg *IdxCfg) SetName(tbl string) *IdxCfg {
	cfg.Index = tbl
	return cfg
}

func (cfg *IdxCfg) SetMapping(m doc.Mapping) *IdxCfg {
	cfg.Mapping = m
	return cfg
}

func (cfg *IdxCfg) SetCustomID(id string) *IdxCfg {
	cfg.UID = id
	return cfg
}

func NewCfgTbl(tbl string, m doc.Mapping, id string) *IdxCfg {
	return NewCfg().
		SetMapping(m).
		SetCustomID(id).
		SetName(tbl)
}

func DefaultCfg() *IdxCfg {
	return NewCfg().
		SetMapping(doc.DefaultMapping()).
		SetName(defaultTbl)
}

func NewCfgFromParams(settings string) (*IdxCfg, error) {
	params, err := param.Parse(settings)
	if err != nil {
		return nil, err
	}
	cfg := NewCfgTbl(params.Index, NewMappingFromParams(params), params.UID)
	return cfg, nil
}

func NewMappingFromParamCfg(cfg *param.Cfg) doc.Mapping {
	if !cfg.HasSrchAttr() && !cfg.HasFacetAttr() {
		return doc.DefaultMapping()
	}

	m := doc.NewMapping()

	for _, attr := range cfg.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range cfg.FacetAttr {
		m.AddKeywords(attr)
	}

	return m
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

func (tbl *IdxCfg) WithCustomID(name string) *IdxCfg {
	tbl.UID = name
	return tbl
}

func (c *IdxCfg) SetID(id int) {
	c.ID = id
}

func (c *IdxCfg) GetID() int {
	return c.ID
}

func (c *IdxCfg) AfterFind(db *hare.Database) error {
	//println("after find")
	return nil
}
