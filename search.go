package srch

import (
	"net/url"
	"strconv"
)

type Search struct {
	*Index
	params url.Values
}

func NewSearch(idx *Index, params string) *Search {
	return &Search{
		Index:  idx,
		params: ParseQuery(params),
	}
}

func ParseRequest(req string) *Search {
	return &Search{
		params: ParseQuery(req),
	}
}

func (r Search) Query() string {
	return r.params.Get(ParamQuery)
}

func (r Search) Filters() *Filters {
	f, _ := DecodeFilter(r.params.Get(ParamFilters))
	return f
}

func (r Search) FacetFilters() *Filters {
	f, _ := DecodeFilter(r.params.Get(ParamFilters))
	return f
}

func (r Search) Facets() []string {
	if r.params.Has(ParamFacets) {
		return r.params[ParamFacets]
	}
	return []string{}
}

func (r Search) Page() int {
	p := r.params.Get(Page)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}

func (r Search) HitsPerPage() int {
	p := r.params.Get(HitsPerPage)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}
