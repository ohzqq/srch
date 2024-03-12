package param

import (
	"encoding/json"
	"net/url"

	"github.com/spf13/cast"
)

type Search struct {
	Hits                 int      `query:"hits,omitempty" json:"hits,omitempty" mapstructure:"hits,omitempty"`
	AttributesToRetrieve []string `query:"attributesToRetrieve,omitempty" json:"attributesToRetrieve,omitempty" mapstructure:"attributesToRetrieve,omitempty"`
	Page                 int      `query:"page,omitempty" json:"page,omitempty" mapstructure:"page,omitempty"`
	HitsPerPage          int      `query:"hitsPerPage,omitempty" json:"hitsPerPage,omitempty" mapstructure:"hitsPerPage,omitempty"`
	SortFacetsBy         string   `query:"sortFacetsBy,omitempty" json:"sortFacetsBy,omitempty" mapstructure:"sortFacetsBy,omitempty"`
	Query                string   `query:"query,omitempty" json:"query,omitempty" mapstructure:"query,omitempty"`
	Facets               []string `query:"facets,omitempty" json:"facets,omitempty" mapstructure:"facets,omitempty"`
	Filters              string   `query:"filters,omitempty" json:"filters,omitempty" mapstructure:"filters,omitempty"`
	FacetFilters         []any    `query:"facetFilters,omitempty" json:"facetFilters,omitempty" mapstructure:"facetFilters,omitempty"`
	NbHits               int      `query:"nbHits,omitempty" json:"nbHits,omitempty" mapstructure:"nbHits,omitempty"`
	NbPages              int      `query:"nbPages,omitempty" json:"nbPages,omitempty" mapstructure:"nbPages,omitempty"`
	SortBy               string   `query:"sortBy,omitempty" json:"sortBy,omitempty" mapstructure:"sortBy,omitempty"`
	Order                string   `query:"order,omitempty" json:"order,omitempty" mapstructure:"order,omitempty"`

	params url.Values
}

func NewSearch() *Search {
	return &Search{
		params: make(url.Values),
	}
}

func (s *Search) Parse(v url.Values) error {
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
		case NbHits:
			s.NbHits = GetQueryInt(key, v)
		case NbPages:
			s.NbPages = GetQueryInt(key, v)
		case SortBy:
			s.SortBy = v.Get(key)
		case Order:
			s.Order = v.Get(key)
		}
		v.Del(key)
	}
	return nil
}

func (p Search) HasFilters() bool {
	return p.params.Has(FacetFilters)
}

//func (p *Search) Filters() []any {
//  if p.HasFilters() {
//    fils, err := unmarshalFilter(p.params.Get(FacetFilters))
//    if err != nil {
//    }
//    return fils
//  }
//  return []any{}
//}

//func (p Search) HitsPerPage() int {
//  page := viper.GetInt(HitsPerPage)
//  if p.params.Has(HitsPerPage) {
//    pn := p.params.Get(HitsPerPage)
//    page, err := strconv.Atoi(pn)
//    if err != nil {
//      return 25
//    }
//    return page
//  }
//  return page
//}

func (p Search) SetHitsPerPage(i any) {
	p.params.Set(HitsPerPage, cast.ToString(i))
}

//func (p Search) Query() string {
//  return p.params.Get(Query)
//}

//func (p Search) SortBy() string {
//  if p.params.Has(SortBy) {
//    return p.params.Get(SortBy)
//  }
//  return DefaultField
//}

//func (p *Search) SortFacetsBy() string {
//  sort := SortByCount
//  if p.params.Has(SortFacetsBy) {
//    if by := p.params.Get(SortFacetsBy); by == SortByCount || by == SortByAlpha {
//      sort = by
//    }
//  }
//  return sort
//}

//func (p Search) Page() int {
//  pn := p.params.Get(Page)
//  page, err := strconv.Atoi(pn)
//  if err != nil {
//    return 0
//  }
//  return page
//}

func (p *Search) SetPage(i any) {
	p.params.Set(Page, cast.ToString(i))
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
