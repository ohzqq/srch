package srch

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ohzqq/sp"
	"github.com/samber/lo"
)

type Cfg struct {
	Search *Search `qs:"-"`
	Idx    *Idx    `qs:"-"`

	//data locations
	Hare   string `json:"-" mapstructure:"db" qs:"db"`
	Data   string `json:"-" mapstructure:"data" qs:"data"`
	IdxURL string `json:"-" mapstructure:"idx" qs:"idx"`
}

func newCfg() *Cfg {
	return &Cfg{
		Idx:    NewIdxCfg(),
		Search: NewSearch(),
	}
}

func NewCfg(v url.Values) (*Cfg, error) {
	cfg := newCfg()
	err := cfg.Decode(v)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Cfg) Decode(v url.Values) error {
	err := Decode(v, cfg.Idx)
	if err != nil {
		return err
	}
	err = Decode(v, cfg.Search)
	if err != nil {
		return err
	}
	err = sp.Decode(v, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) Encode() (url.Values, error) {
	iv, err := cfg.Idx.Encode()
	if err != nil {
		return nil, err
	}
	sv, err := cfg.Search.Encode()
	if err != nil {
		return nil, err
	}
	cv, err := sp.Encode(cfg)
	if err != nil {
		return nil, err
	}
	return lo.Assign(iv, sv, cv), nil
}

func (cfg *Cfg) SetIdxName(tbl string) *Cfg {
	cfg.Idx.Name = tbl
	return cfg
}

func (cfg *Cfg) SetCustomID(id string) *Cfg {
	cfg.Idx.PrimaryKey = id
	return cfg
}

func (cfg *Cfg) IndexName() string {
	if cfg.Idx.Name != "" {
		return cfg.Idx.Name
	}
	return "default"
}

func (cfg *Cfg) DB() *url.URL {
	u, err := parseURL(cfg.Hare)
	if err != nil {
		return &url.URL{Scheme: "mem"}
	}
	return u
}

func (cfg *Cfg) SrchURL() *url.URL {
	u, err := parseURL(cfg.IdxURL)
	if err != nil {
		return &url.URL{Scheme: "mem"}
	}
	return u
}

func (cfg *Cfg) DataURL() *url.URL {
	u, err := parseURL(cfg.Data)
	if err != nil {
		u = &url.URL{Scheme: "mem"}
	}
	if cfg.Idx.PrimaryKey != "" {
		v := u.Query()
		v.Set("primaryKey", cfg.Idx.PrimaryKey)
		u.RawQuery = v.Encode()
	}
	return u
}

func (cfg *Cfg) HasData() bool {
	return cfg.Data != ""
}

func (cfg *Cfg) HasDB() bool {
	return cfg.Hare != ""
}

func (cfg *Cfg) HasIdxURL() bool {
	return cfg.IdxURL != ""
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

func parseURL(uri string) (*url.URL, error) {
	var err error
	u := &url.URL{}
	if uri == "" {
		return u, nil
	}
	u, err = url.Parse(uri)
	if err != nil {
		return u, err
	}
	if u.Scheme == "file" {
		u.Path = filepath.Join("/", u.Host, u.Path)
		//u.Host = ""
	}
	return u, nil
}

func parseSrchAttrs(attrs []string) []string {
	switch len(attrs) {
	case 0:
		return []string{"*"}
	case 1:
		if attrs[0] == "" {
			return []string{"*"}
		}
		fallthrough
	default:
		return ParseQueryStrings(attrs)
	}
}

func ParseQueryStrings(q []string) []string {
	var vals []string
	for _, val := range q {
		if val == "" {
			break
		}
		for _, v := range strings.Split(val, ",") {
			vals = append(vals, v)
		}
	}
	return vals
}
