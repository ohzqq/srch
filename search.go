package srch

import (
	"net/url"
	"strconv"
)

const (
	SearchQuery   = "query"
	SearchFilters = "filters"
	FacetFilters  = "facetFilters"
	SearchFacets  = "facets"
	Page          = "page"
	HitsPerPage   = "hitsPerPage"
)

type Request struct {
	params url.Values
}

func ParseRequest(req string) *Request {
	return &Request{
		params: NewQuery(req),
	}
}

func (r Request) Query() string {
	return r.params.Get(SearchQuery)
}

func (r Request) Filters() *Filters {
	f, _ := DecodeFilter(r.params.Get(SearchFilters))
	return f
}

func (r Request) FacetFilters() *Filters {
	f, _ := DecodeFilter(r.params.Get(FacetFilters))
	return f
}

func (r Request) Facets() []string {
	if r.params.Has(SearchFacets) {
		return r.params[SearchFacets]
	}
	return []string{}
}

func (r Request) Page() int {
	p := r.params.Get(Page)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}

func (r Request) HitsPerPage() int {
	p := r.params.Get(HitsPerPage)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}
