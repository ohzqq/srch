package param

import (
	"encoding/json"
	"net/url"
)

type Cfg struct {
	Client *Client
	Search *Search
	Idx    *Idx
}

func NewCfg() *Cfg {
	return &Cfg{
		Idx:    NewIdx(),
		Search: NewSearch(),
		Client: NewClient(),
	}
}

func (cfg *Cfg) Decode(v any) error {
	err := Decode(v, cfg.Idx)
	if err != nil {
		return err
	}
	err = Decode(v, cfg.Search)
	if err != nil {
		return err
	}
	err = Decode(v, cfg.Client)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Cfg) IndexName() string {
	if cfg.Client.Index != "" {
		return cfg.Client.Index
	}
	return "default"
}

func (cfg *Cfg) DB() *url.URL {
	u, err := parseURL(cfg.Client.DB)
	if err != nil {
		return &url.URL{Scheme: "mem"}
	}
	return u
}

func (cfg *Cfg) SrchURL() *url.URL {
	u, err := parseURL(cfg.Search.URI)
	if err != nil {
		return &url.URL{Scheme: "mem"}
	}
	return u
}

func (cfg *Cfg) DataURL() *url.URL {
	u, err := parseURL(cfg.Idx.Data)
	if err != nil {
		return &url.URL{Scheme: "mem"}
	}
	return u
}

func (cfg *Cfg) HasData() bool {
	return cfg.Idx.Data != ""
}

func (cfg *Cfg) HasDB() bool {
	return cfg.Client.DB != ""
}
func (cfg *Cfg) HasIdx() bool {
	return cfg.Search.URI != ""
}

func (cfg *Cfg) HasSrchAttr() bool {
	return len(cfg.Idx.SrchAttr) > 0
}

func (cfg *Cfg) HasFacetAttr() bool {
	return len(cfg.Idx.FacetAttr) > 0
}

func (cfg *Cfg) HasSortAttr() bool {
	return len(cfg.Idx.SortAttr) > 0
}

func (cfg *Cfg) HasFilters() bool {
	return len(cfg.Search.FacetFltr) > 0
}

func (cfg *Cfg) Filters() []any {
	if len(cfg.Search.FacetFltr) > 0 {
		var fltr []any
		err := json.Unmarshal([]byte(cfg.Search.FacetFltr[0]), &fltr)
		if err != nil {
			return []any{""}
		}
		return fltr
	}
	return []any{""}
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
