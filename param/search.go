package param

import (
	"encoding/json"
	"net/url"

	"github.com/spf13/viper"
)

type Search struct {
	Hits                 int      `query:"hits,omitempty" json:"hits,omitempty"`
	AttributesToRetrieve []string `query:"attributesToRetrieve,omitempty" json:"attributesToRetrieve,omitempty"`
	Page                 int      `query:"page,omitempty" json:"page,omitempty"`
	HitsPerPage          int      `query:"hitsPerPage,omitempty" json:"hitsPerPage,omitempty"`
	SortFacetsBy         string   `query:"sortFacetsBy,omitempty" json:"sortFacetsBy,omitempty"`
	Query                string   `query:"query,omitempty" json:"query,omitempty"`
	Facets               []string `query:"facets,omitempty" json:"facets,omitempty"`
	Filters              string   `query:"filters,omitempty" json:"filters,omitempty"`
	FacetFilters         []any    `query:"facetFilters,omitempty" json:"facetFilters,omitempty"`
	SortBy               string   `query:"sortBy,omitempty" json:"sortBy,omitempty"`
	Order                string   `query:"order,omitempty" json:"order,omitempty"`
}

func NewSearch() *Search {
	return &Search{
		HitsPerPage: viper.GetInt("hitsPerPage"),
	}
}

func (s *Search) Parse(q string) error {
	vals, err := url.ParseQuery(q)
	if err != nil {
		return err
	}
	return s.Set(vals)
}

func (s *Search) Set(v url.Values) error {
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
		case SortFacetsBy:
			s.SortFacetsBy = v.Get(key)
		case Query:
			s.Query = v.Get(key)
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
		case SortBy:
			s.SortBy = v.Get(key)
		case Order:
			s.Order = v.Get(key)
		}
		v.Del(key)
	}
	return nil
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
