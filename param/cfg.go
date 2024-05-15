package param

import (
	"net/url"

	"github.com/ohzqq/sp"
)

type Cfg struct {
	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes"`
	ID        string   `query:"id,omitempty" json:"id,omitempty" mapstructure:"id" qs:"id"`
	Index     string   `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	Path      string   `json:"-" mapstructure:"path" qs:"path"`
}

func NewCfg() *Cfg {
	return &Cfg{}
}

func ParseCfg(q string) (*Cfg, error) {
	cfg := NewCfg()
	err := cfg.Decode(q)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Cfg) Decode(q string) error {
	u, err := url.Parse(q)
	if err != nil {
		return err
	}
	err = sp.Decode(u.Query(), cfg)
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

func (cfg *Cfg) Encode() (url.Values, error) {
	return sp.Encode(cfg)
}
