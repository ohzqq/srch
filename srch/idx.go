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

func (idx *Idx) SetMapping(m Mapping) *Idx {
	idx.Mapping = m
	return idx
}

func (idx *Idx) Decode(u url.Values) error {
	err := sp.Decode(u, idx)
	if err != nil {
		return err
	}
	idx.SrchAttr = parseSrchAttrs(idx.SrchAttr)
	if len(idx.FacetAttr) > 0 {
		idx.FacetAttr = ParseQueryStrings(idx.FacetAttr)
	}
	if len(idx.SortAttr) > 0 {
		idx.SortAttr = ParseQueryStrings(idx.SortAttr)
	}
	idx.SetMapping(idx.mapParams())
	return nil
}

func (idx *Idx) HasData() bool {
	return idx.Data != ""
}

func (idx *Idx) HasSrchAttr() bool {
	return len(idx.SrchAttr) > 0
}

func (idx *Idx) HasFacetAttr() bool {
	return len(idx.FacetAttr) > 0
}

func (idx *Idx) HasSortAttr() bool {
	return len(idx.SortAttr) > 0
}

func (idx *Idx) Encode() (url.Values, error) {
	return sp.Encode(idx)
}

func (idx *Idx) mapParams() Mapping {
	m := NewMapping()

	for _, attr := range idx.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range idx.FacetAttr {
		m.AddKeywords(attr)
	}

	for _, attr := range idx.SortAttr {
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
