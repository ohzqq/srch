package srch

import (
	"net/url"
	"strconv"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
)

type Request struct {
	*Index
	params url.Values
}

func Search(idx *Index, params string) *Response {
	req := ParseRequest(params)

	q := req.Query()
	if q == "" {
		return NewResponse(idx.Data, req.params)
	}

	sids := idx.FuzzySearch(q)
	res := roaring.And(idx.Clone(), sids)
	f := FilteredItems(idx.Data, lo.ToAnySlice(res.ToArray()))
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
