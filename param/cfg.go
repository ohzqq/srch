package param

import (
	"net/url"

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
	if cfg.Path != "" {
		cfg.URL, err = url.Parse(cfg.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cfg *Cfg) Encode() (url.Values, error) {
	return sp.Encode(cfg)
}
