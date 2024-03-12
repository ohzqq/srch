package param

import (
	"net/url"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Params struct {
	// Settings
	*Settings `mapstructure:",squash"`
	*Search   `mapstructure:",squash"`
	Other     url.Values `mapstructure:"params,omitempty"`
}

func New() *Params {
	return &Params{
		Settings: NewSettings(),
		Search:   NewSearch(),
		Other:    make(url.Values),
	}
}

func Parse(params string) (*Params, error) {
	p := New()

	vals, err := url.ParseQuery(params)
	if err != nil {
		return nil, err
	}

	err = p.Settings.Parse(vals)
	if err != nil {
		return nil, err
	}
	err = p.Search.Parse(vals)
	if err != nil {
		return nil, err
	}

	p.Other = vals

	return p, nil
}

// ParseParams takes an interface{} and returns a url.Values.
func ParseParams(f string) (url.Values, error) {
	vals, err := url.ParseQuery(f)
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
	if !vals.Has(Facets) {
		vals[Facets] = GetQueryStringSlice(FacetAttr, vals)
	}
	return vals[Facets]
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, err
}

func GetAnySlice(key string, vals url.Values) []any {
	return lo.ToAnySlice(GetQueryStringSlice(key, vals))
}

func GetQueryInt(key string, vals url.Values) int {
	if vals.Has(key) {
		return cast.ToInt(vals.Get(key))
	}
	return 0
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
