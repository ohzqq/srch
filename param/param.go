package param

import (
	"net/url"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Params struct {
	// Settings
	*IndexSettings `mapstructure:",squash"`
	*Search        `mapstructure:",squash"`
	*SrchCfg
	params url.Values
	Other  url.Values `mapstructure:"params,omitempty"`
}

func New() *Params {
	return &Params{
		IndexSettings: NewSettings(),
		Search:        NewSearch(),
		SrchCfg:       NewCfg(),
		Other:         make(url.Values),
		params:        make(url.Values),
	}
}

func Parse(params string) (*Params, error) {
	p := New()

	vals, err := url.ParseQuery(params)
	if err != nil {
		return nil, err
	}
	p.params = vals

	//err = p.SrchCfg.Set(vals)
	//if err != nil {
	//return nil, err
	//}
	//err = p.IndexSettings.Set(vals)
	//if err != nil {
	//return nil, err
	//}
	//err = p.Search.Set(vals)
	//if err != nil {
	//return nil, err
	//}
	err = p.Set(vals)
	if err != nil {
		return nil, err
	}

	p.Other = vals

	return p, nil
}

func (s *Params) Set(v url.Values) error {
	for _, key := range paramsSettings {
		switch key {
		case SrchAttr:
			s.SrchAttr = parseSrchAttr(v)
		case FacetAttr:
			s.FacetAttr = parseFacetAttr(v)
		case SortAttr:
			s.SortAttr = GetQueryStringSlice(key, v)
		case DefaultField:
			s.DefaultField = v.Get(key)
		case UID:
			s.IndexSettings.UID = v.Get(key)
		}
	}
	for _, key := range paramsCfg {
		switch key {
		case DataDir:
			s.DataDir = v.Get(key)
		case DataFile:
			s.DataFile = GetQueryStringSlice(key, v)
		case FullText:
			s.BlvPath = v.Get(key)
		case UID:
			s.SrchCfg.UID = v.Get(key)
		}
	}
	for _, key := range paramsFacets {
		switch key {
		case SortFacetsBy:
			s.SortFacetsBy = v.Get(key)
		case Facets:
			s.Facets = GetQueryStringSlice(key, v)
		case Filters:
			s.Filters = v.Get(key)
		case FacetFilters:
			if v.Has(key) {
				fil := v.Get(key)
				f, err := unmarshalFilter(fil)
				if err != nil {
					return err
				}
				s.FacetFilters = f
			}
		}
	}
	for _, key := range paramsSearch {
		switch key {
		case Hits:
			s.Hits = GetQueryInt(key, v)
		case AttributesToRetrieve:
			s.AttributesToRetrieve = GetQueryStringSlice(key, v)
		case Page:
			s.Page = GetQueryInt(key, v)
		case HitsPerPage:
			s.HitsPerPage = GetQueryInt(key, v)
		case Query:
			s.Query = v.Get(key)
		case SortBy:
			s.SortBy = v.Get(key)
		case Order:
			s.Order = v.Get(key)
		}
	}
	return nil
}

func (s *Params) Has(key string) bool {
	switch key {
	case Hits:
		return s.Hits != 0
	case AttributesToRetrieve:
		return len(s.AttributesToRetrieve) != 0
	case Page:
		return s.Page != 0
	case HitsPerPage:
		return s.HitsPerPage != 0
	case Query:
		return s.Query != ""
	case SortBy:
		return s.SortBy != ""
	case Order:
		return s.Order != ""
	case DataDir:
		return s.DataDir != ""
	case DataFile:
		return len(s.DataFile) > 0
	case FullText:
		return s.BlvPath != ""
	case UID:
		return s.IndexSettings.UID != "" || s.SrchCfg.UID != ""
	case SortFacetsBy:
		return s.SortFacetsBy != ""
	case Facets:
		return len(s.Facets) > 0
	case Filters:
		return s.Filters != ""
	case FacetFilters:
		return len(s.FacetFilters) > 0
	case SrchAttr:
		return len(s.SrchAttr) > 0
	case FacetAttr:
		return len(s.FacetAttr) > 0
	case SortAttr:
		return len(s.SortAttr) > 0
	case DefaultField:
		return s.DefaultField != ""
	default:
		return false
	}
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
