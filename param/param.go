package param

import (
	"net/url"
	"strings"

	"github.com/spf13/cast"
)

type Params struct {
	*Settings
	*Search
	params url.Values
}

func New() *Params {
	return &Params{
		Settings: NewSettings(),
		Search:   NewSearch(),
		params:   make(url.Values),
	}
}

func Parse(q string) *Params {
}

func (p *Params) Get(key string) string {
	if p.params.Has(key) {
		return p.params.Get(key)
	}
	if p.Search.params.Has(key) {
		return p.Search.params.Get(key)
	}
}

func (p *Params) All(key string) []string {
	if p.Settings.params.Has(key) {
		return p.Settings.params[key]
	}
	if p.Search.params.Has(key) {
		return p.Search.params[key]
	}
}

func (p *Params) Del(key string) {
	if p.Settings.params.Has(key) {
		p.Settings.params.Del(key)
	}
	if p.Search.params.Has(key) {
		p.Search.params.Del(key)
	}
}

func (p *Params) Set(key string, val any) {
	v := cast.ToString(val)
	if p.Settings.params.Has(key) {
		p.Settings.params.Set(key, v)
	}
	if p.Search.params.Has(key) {
		p.Search.params.Set(key, v)
	}
}

func (p *Params) Add(key string, val any) {
	v := cast.ToString(val)
	if p.Settings.params.Has(key) {
		p.Settings.params.Add(key, v)
	}
	if p.Search.params.Has(key) {
		p.Search.params.Add(key, v)
	}
}

// ParseParams takes an interface{} and returns a url.Values.
func ParseParams(f string) (url.Values, error) {
	vals, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}

	vals[SrchAttr] = parseSrchAttr(vals)
	vals[FacetAttr] = parseFacetAttr(vals)
	return vals, nil
}

func parseSrchAttr(vals url.Values) []string {
	if !vals.Has(SrchAttr) {
		return []string{DefaultField}
	}
	vals[SrchAttr] = GetQueryStringSlice(SrchAttr, vals)
	if len(vals[SrchAttr]) < 1 {
		vals[SrchAttr] = []string{DefaultField}
	}
	return vals[SrchAttr]
}

func parseFacetAttr(vals url.Values) []string {
	if !vals.Has(ParamFacets) {
		vals[ParamFacets] = GetQueryStringSlice(FacetAttr, vals)
	}
	return vals[ParamFacets]
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, err
}

func GetQueryStringSlice(key string, q url.Values) []string {
	var vals []string
	if q.Has(key) {
		for _, val := range q[key] {
			if val == "" {
				break
			}
			for _, v := range strings.Split(val, ",") {
				vals = append(vals, v)
			}
		}
	}
	return vals
}
