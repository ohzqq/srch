package param

import "net/url"

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

func (cfg *Cfg) HasSrchAttr() bool {
	return len(cfg.Idx.SrchAttr) > 0
}

func (cfg *Cfg) HasFacetAttr() bool {
	return len(cfg.Idx.FacetAttr) > 0
}

func (cfg *Cfg) HasSortAttr() bool {
	return len(cfg.Idx.SortAttr) > 0
}
