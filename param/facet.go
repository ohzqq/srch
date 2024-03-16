package param

import (
	"encoding/json"
	"net/url"
)

type FacetSettings struct {
	UID          string   `query:"uid,omitempty" json:"uid,omitempty"`
	Facets       []string `query:"facets,omitempty" json:"facets,omitempty"`
	Filters      string   `query:"filters,omitempty" json:"filters,omitempty"`
	FacetFilters []any    `query:"facetFilters,omitempty" json:"facetFilters,omitempty"`
	SortFacetsBy string   `query:"sortFacetsBy,omitempty" json:"sortFacetsBy,omitempty"`
}

func NewFacetSettings() *FacetSettings {
	return &FacetSettings{}
}

func (facet *FacetSettings) Set(v url.Values) error {
	for _, key := range paramsFacets {
		switch key {
		case SortFacetsBy:
			facet.SortFacetsBy = v.Get(key)
		case Facets:
			facet.Facets = GetQueryStringSlice(key, v)
		case Filters:
			facet.Filters = v.Get(key)
		case FacetFilters:
			if v.Has(key) {
				fil := v.Get(key)
				f, err := unmarshalFilter(fil)
				if err != nil {
					return err
				}
				facet.FacetFilters = f
			}
		}
		v.Del(key)
	}
	return nil
}

func (s *FacetSettings) Has(key string) bool {
	switch key {
	case SortFacetsBy:
		return s.SortFacetsBy != ""
	case Facets:
		return len(s.Facets) > 0
	case Filters:
		return s.Filters != ""
	case FacetFilters:
		return len(s.FacetFilters) > 0
	}
	return false
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
