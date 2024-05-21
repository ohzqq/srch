package srch

import (
	"net/url"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
)

type IdxCfg struct {
	ID      int     `json:"_id"`
	Mapping Mapping `json:"mapping"`
	Name    string  `json:"name" qs:"index"`

	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes,omitempty"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting,omitempty"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes,omitempty"`
	Data      string   `json:"-" mapstructure:"path" qs:"data"`
}

func NewIdxCfg() *IdxCfg {
	cfg := &IdxCfg{
		SrchAttr: []string{"*"},
		Name:     "default",
	}
	return cfg.
		SetMapping(DefaultMapping())
}

func (cfg *IdxCfg) SetMapping(m Mapping) *IdxCfg {
	cfg.Mapping = m
	return cfg
}

func (cfg *IdxCfg) Decode(u url.Values) error {
	err := sp.Decode(u, cfg)
	if err != nil {
		return err
	}
	cfg.SrchAttr = parseSrchAttrs(cfg.SrchAttr)
	if len(cfg.FacetAttr) > 0 {
		cfg.FacetAttr = ParseQueryStrings(cfg.FacetAttr)
	}
	if len(cfg.SortAttr) > 0 {
		cfg.SortAttr = ParseQueryStrings(cfg.SortAttr)
	}
	cfg.SetMapping(cfg.mapParams())
	return nil
}

func (cfg *IdxCfg) Encode() (url.Values, error) {
	return sp.Encode(cfg)
}

func (cfg *IdxCfg) mapParams() Mapping {
	m := NewMapping()

	for _, attr := range cfg.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range cfg.FacetAttr {
		m.AddKeywords(attr)
	}

	for _, attr := range cfg.SortAttr {
		m.AddKeywords(attr)
	}

	return m
}

//func CfgEqual(old, cur *Idx) bool {
//  if !slices.Equal(old.SrchAttr, cur.SrchAttr) {
//    return false
//  }
//  if !slices.Equal(old.FacetAttr, cur.FacetAttr) {
//    return false
//  }
//  if !slices.Equal(old.SortAttr, cur.SortAttr) {
//    return false
//  }
//  if old.Index != cur.Index {
//    return false
//  }
//  if old.UID != cur.UID {
//    return false
//  }
//  return true
//}

func NewCfgTbl(tbl string, m Mapping, id string) *IdxCfg {
	return NewIdxCfg().
		SetMapping(m)
}

func DefaultIdxCfg() *IdxCfg {
	return NewIdxCfg().
		SetMapping(DefaultMapping())
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
