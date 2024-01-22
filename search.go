package srch

import (
	"log"
	"net/url"
	"strconv"

	"github.com/samber/lo"
)

type Request struct {
	*Index
	params url.Values
}

func Search(idx *Index, params string) *Response {
	req := ParseRequest(params)

	var filtered []int
	if req.params.Has(ParamFacetFilters) {
		var err error
		filtered, err = Filter(idx.Fields, req.params.Get(ParamFacetFilters))
		if err != nil {
			log.Fatal(err)
		}
	}

	q := req.Query()
	if q == "" {
		f := FilteredItems(idx.Data, lo.ToAnySlice(filtered))
		return NewResponse(f, req.params)
	}

	sids := idx.FuzzySearch(q)
	if len(filtered) > 0 {
		sids = lo.Intersect(sids, filtered)
	}
	f := FilteredItems(idx.Data, lo.ToAnySlice(sids))
	return NewResponse(f, req.params)
}

func ParseRequest(req string) *Request {
	r := &Request{
		params: ParseQuery(req),
	}
	return r
}

func (r Request) Query() string {
	return r.params.Get(ParamQuery)
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
