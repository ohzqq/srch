package param

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/ohzqq/sp"
)

type Search struct {
	// Search
	Hits        int      `query:"hits,omitempty" json:"hits,omitempty" mapstructure:"hits" qs:"hits"`
	RtrvAttr    []string `query:"attributesToRetrieve,omitempty" json:"attributesToRetrieve,omitempty" mapstructure:"attributes_to_retrieve" qs:"attributesToRetrieve"`
	Page        int      `query:"page,omitempty" json:"page,omitempty" mapstructure:"page" qs:"page"`
	HitsPerPage int      `query:"hitsPerPage,omitempty" json:"hitsPerPage,omitempty" mapstructure:"hits_per_page" qs:"hitsPerPage"`
	Query       string   `query:"query,omitempty" json:"query,omitempty" mapstructure:"query" qs:"query"`
	SortBy      string   `query:"sortBy,omitempty" json:"sortBy,omitempty" mapstructure:"sort_by" qs:"sortBy"`
	Order       string   `query:"order,omitempty" json:"order,omitempty" mapstructure:"order" qs:"order"`

	// Facets
	Facets       []string `query:"facets,omitempty" json:"facets,omitempty" mapstructure:"facets" qs:"facets"`
	Filters      string   `query:"filters,omitempty" json:"filters,omitempty" mapstructure:"filters" qs:"filters"`
	FacetFltr    []string `query:"facetFilters,omitempty" json:"facetFilters,omitempty" mapstructure:"facet_filters" qs:"facetFilters"`
	SortFacetsBy string   `query:"sortFacetsBy,omitempty" json:"sortFacetsBy,omitempty" mapstructure:"sort_facets_by" qs:"sortFacetsBy"`
	MaxFacetVals int      `query:"maxValuesPerFacet,omitempty" json:"maxValuesPerFacet,omitempty" mapstructure:"max_values_per_facet" qs:"maxValuesPerFacet,omitempty"`

	ID    string `query:"id,omitempty" json:"id,omitempty" mapstructure:"id" qs:"id"`
	Index string `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	Path  string `json:"-" mapstructure:"path" qs:"path"`
}

func NewSearch() *Search {
	return &Search{
		Index: "default",
	}
}

func ParseSearch(q string) (*Search, error) {
	s := NewSearch()
	err := s.Decode(q)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Search) Decode(q string) error {
	q = strings.TrimPrefix(q, "?")
	u, err := url.ParseQuery(q)
	if err != nil {
		return err
	}
	err = sp.Decode(u, s)
	if err != nil {
		return err
	}
	s.RtrvAttr = parseSrchAttrs(s.RtrvAttr)
	s.Facets = ParseQueryStrings(s.Facets)
	return nil
}

func (s *Search) Encode() (url.Values, error) {
	v, err := sp.Encode(s)
	if err != nil {
		return nil, err
	}
	if fltr := s.FacetFilters(); len(fltr) > 0 {
		d, err := json.Marshal(fltr)
		if err != nil {
			return nil, err
		}
		v.Set(FacetFilters.String(), string(d))
	}
	return v, nil
}

func (s *Search) FacetFilters() []any {
	if len(s.FacetFltr) > 0 {
		var fltr []any
		err := json.Unmarshal([]byte(s.FacetFltr[0]), &fltr)
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
