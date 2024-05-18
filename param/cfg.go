package param

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
	URI       string   `json:"-" mapstructure:"path" qs:"url"`
}

func NewIdx() *Idx {
	return &Idx{
		Client:   NewClient(),
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
	cfg.URL, err = parseURL(cfg.URI)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Idx) Encode() (url.Values, error) {
	return sp.Encode(cfg)
}

func (cfg *Idx) HasSrchAttr() bool {
	return len(cfg.SrchAttr) > 0
}

func (cfg *Idx) HasFacetAttr() bool {
	return len(cfg.FacetAttr) > 0
}

func (cfg *Idx) HasSortAttr() bool {
	return len(cfg.SortAttr) > 0
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
