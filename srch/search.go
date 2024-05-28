package srch

import (
	"net/url"

	"github.com/ohzqq/sp"
	"github.com/samber/lo"
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

	*url.URL `json:"-"`
}

func NewSearch() *Search {
	return &Search{}
}

func (s *Search) FilterRtrvAttr(data []map[string]any) []map[string]any {
	if len(data) < 1 {
		return data
	}
	for _, attr := range s.RtrvAttr {
		if attr == "*" {
			return data
		}
		fn := func(d map[string]any, _ int) map[string]any {
			return lo.PickByKeys(d, s.RtrvAttr)
		}
		return lo.Map(data, fn)
	}
	return data
}

func (s *Search) Decode(u url.Values) error {
	err := sp.Decode(u, s)
	if err != nil {
		return err
	}
	s.RtrvAttr = parseSrchAttrs(s.RtrvAttr)
	if len(s.Facets) > 0 {
		s.Facets = ParseQueryStrings(s.Facets)
	}
	return nil
}

func (s *Search) Encode() (url.Values, error) {
	v, err := sp.Encode(s)
	if err != nil {
		return nil, err
	}
	//if fltr := s.FacetFilters(); len(fltr) > 0 {
	//  d, err := json.Marshal(fltr)
	//  if err != nil {
	//    return nil, err
	//  }
	//  v.Set(FacetFilters.String(), string(d))
	//}

	return v, nil
}
