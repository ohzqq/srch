package param

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Search struct {
	params url.Values
}

func NewSearch() *Search {
	return &Search{
		params: make(url.Values),
	}
}

func (p Search) HasFilters() bool {
	return p.params.Has(FacetFilters)
}

func (p *Search) Filters() []any {
	if p.HasFilters() {
		fils, err := unmarshalFilter(p.Get(FacetFilters))
		if err != nil {
		}
		return fils
	}
	return p.filters
}

func (p Search) HitsPerPage() int {
	page := viper.GetInt(HitsPerPage)
	if p.params.Has(HitsPerPage) {
		pn := p.params.Get(HitsPerPage)
		page, err := strconv.Atoi(pn)
		if err != nil {
			return 25
		}
		return page
	}
	return page
}

func (p Search) SetHitsPerPage(i any) {
	p.params.Set(HitsPerPage, cast.ToString(i))
}

func (p Search) Query() string {
	return p.Get(Query)
}

func (p Search) SortBy() string {
	if p.Has(SortBy) {
		return p.Get(SortBy)
	}
	return DefaultField
}

//func (p *Search) SortFacetsBy() string {
//  sort := SortByCount
//  if p.params.Has(SortFacetsBy) {
//    if by := p.params.Get(SortFacetsBy); by == SortByCount || by == SortByAlpha {
//      sort = by
//    }
//  }
//  return sort
//}

func (p Search) Page() int {
	pn := p.params.Get(Page)
	page, err := strconv.Atoi(pn)
	if err != nil {
		return 0
	}
	return page
}

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
