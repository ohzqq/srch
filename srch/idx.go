package srch

import (
	"net/url"
	"slices"

	"github.com/ohzqq/sp"
)

type Idx struct {
	*url.URL `json:"-"`
	*Client

	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes,omitempty"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting,omitempty"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes,omitempty"`
	Data      string   `json:"-" mapstructure:"path" qs:"data"`
}

func NewIdx() *Idx {
	client, _ := NewClient("")
	return &Idx{
		Client:   client,
		SrchAttr: []string{"*"},
	}
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
	return nil
}

func (cfg *Idx) Encode() (url.Values, error) {
	return sp.Encode(cfg)
}

func CfgEqual(old, cur *Idx) bool {
	if !slices.Equal(old.SrchAttr, cur.SrchAttr) {
		return false
	}
	if !slices.Equal(old.FacetAttr, cur.FacetAttr) {
		return false
	}
	if !slices.Equal(old.SortAttr, cur.SortAttr) {
		return false
	}
	if old.Index != cur.Index {
		return false
	}
	if old.UID != cur.UID {
		return false
	}
	return true
}
