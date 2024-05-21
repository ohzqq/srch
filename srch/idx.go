package srch

import (
	"net/url"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
)

type Idx struct {
	ID      int     `json:"_id"`
	Mapping Mapping `json:"mapping"`
	Name    string  `json:"name" qs:"index"`

	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes,omitempty"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting,omitempty"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes,omitempty"`
	Data      string   `json:"-" mapstructure:"path" qs:"data"`
}

func NewIdxCfg() *Idx {
	cfg := &Idx{
		SrchAttr: []string{"*"},
		Name:     "default",
	}
	return cfg.
		SetMapping(DefaultMapping())
}

func (cfg *Idx) SetMapping(m Mapping) *Idx {
	cfg.Mapping = m
	return cfg
}

func (cfg *Idx) Decode(u url.Values) error {
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

func (cfg *Idx) Encode() (url.Values, error) {
	return sp.Encode(cfg)
}

func (cfg *Idx) mapParams() Mapping {
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

func NewCfgTbl(tbl string, m Mapping, id string) *Idx {
	return NewIdxCfg().
		SetMapping(m)
}

func DefaultIdxCfg() *Idx {
	return NewIdxCfg().
		SetMapping(DefaultMapping())
}

func (c *Idx) SetID(id int) {
	c.ID = id
}

func (c *Idx) GetID() int {
	return c.ID
}

func (c *Idx) AfterFind(db *hare.Database) error {
	//println("after find")
	return nil
}
