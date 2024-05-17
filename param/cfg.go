package param

import (
	"net/url"
	"path/filepath"
	"slices"

	"github.com/ohzqq/sp"
)

type Cfg struct {
	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes,omitempty"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting,omitempty"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes,omitempty"`

	*Paramz
}

func NewCfg() *Cfg {
	return &Cfg{
		Paramz:   defaultParams(),
		SrchAttr: []string{"*"},
	}
}

func (cfg *Cfg) Decode(u url.Values) error {
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
	if cfg.URI != "" {
		cfg.URL, err = url.Parse(cfg.URI)
		if err != nil {
			return err
		}
		if cfg.URL.Scheme == "file" {
			cfg.URL.Path = filepath.Join("/", cfg.URL.Host, cfg.URL.Path)
		}
	}
	return nil
}

func (cfg *Cfg) Encode() (url.Values, error) {
	return sp.Encode(cfg)
}

func (cfg *Cfg) HasSrchAttr() bool {
	return len(cfg.SrchAttr) > 0
}

func (cfg *Cfg) HasFacetAttr() bool {
	return len(cfg.FacetAttr) > 0
}

func (cfg *Cfg) HasSortAttr() bool {
	return len(cfg.SortAttr) > 0
}

func CfgEqual(old, cur *Cfg) bool {
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
