package srch

import (
	"strconv"

	"github.com/RoaringBitmap/roaring"
)

type Request struct {
	*Index
	*Query
}

func Search(idx *Index, params string) *roaring.Bitmap {
	req := ParseRequest(params)

	q := req.Keywords()
	if q == "" {
		return roaring.New()
	}

	bits := idx.Bitmap()
	bits.And(idx.FuzzySearch(q))
	return bits
}

func ParseRequest(req string) *Request {
	r := &Request{
		Query: NewQuery(req),
	}
	return r
}

func (r Request) Keywords() string {
	return r.Params.Get(ParamQuery)
}

func (r Request) Facets() []string {
	if r.Params.Has(ParamFacets) {
		return r.Params[ParamFacets]
	}
	return []string{}
}

func (r Request) Page() int {
	p := r.Params.Get(Page)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}

func (r Request) HitsPerPage() int {
	p := r.Params.Get(HitsPerPage)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}
