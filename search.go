package srch

import (
	"net/url"
	"strconv"
)

type Request struct {
	params url.Values
}

func ParseRequest(req string) *Request {
	return &Request{
		params: ParseQuery(req),
	}
}

func (r Request) Query() string {
	return r.params.Get(ParamQuery)
}

func (r Request) Filters() *Filters {
	f, _ := DecodeFilter(r.params.Get(ParamFilters))
	return f
}

func (r Request) FacetFilters() *Filters {
	f, _ := DecodeFilter(r.params.Get(ParamFilters))
	return f
}

func (r Request) Facets() []string {
	if r.params.Has(ParamFacets) {
		return r.params[ParamFacets]
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
