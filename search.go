package srch

import (
	"fmt"
	"net/url"
	"strconv"
)

type Request struct {
	*Index
	params url.Values
	*Filters
}

func Search(idx *Index, params string) *Response {
	req := ParseRequest(params)

	if req.HasFilters() {
		println("filters")
	}

	q := req.Query()
	if q == "" {
		return NewResponse(idx.Data, req.params)
	}
	data := idx.search(q)
	return NewResponse(data, req.params)
}

func ParseRequest(req string) *Request {
	r := &Request{
		params: ParseQuery(req),
	}
	fmt.Printf("filters?? %v\n", r.params.Get(ParamFacetFilters))
	r.Filters = r.GetFilters()
	return r
}

func (r Request) Query() string {
	return r.params.Get(ParamQuery)
}

func (r Request) GetFilters() *Filters {
	var f *Filters
	switch {
	case r.params.Has(ParamFilters):
		f, _ = DecodeFilter(r.params.Get(ParamFilters))
	case r.params.Has(ParamFacetFilters):
		return r.getFacetFilters()
	}
	return f
}

func (r Request) HasFilters() bool {
	if r.Filters == nil {
		return false
	}
	return len(r.Filters.Not) > 0 ||
		len(r.Filters.And) > 0 ||
		len(r.Filters.Or) > 0
}

func (r Request) getFacetFilters() *Filters {
	f, _ := DecodeFilter(r.params.Get(ParamFacetFilters))
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
