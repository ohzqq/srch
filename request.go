package srch

import (
	"net/url"

	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
)

type Request struct {
	*param.Params
	filters []string
}

func NewRequest() *Request {
	return &Request{
		Params: param.New(),
	}
}

func ParseRequest(params string) (*Request, error) {
	p, err := param.Parse(params)
	if err != nil {
		return nil, err
	}
	return &Request{Params: p}, nil
}

func (r *Request) Parse(params string) (*Request, error) {
	return ParseRequest(params)
}

func (r *Request) SetValues(vals url.Values) (*Request, error) {
	err := r.Params.Set(vals)
	return r, err
}

func (p *Request) SetRoute(path string) *Request {
	p.URL.Path = path
	p.ParseRoute(p.URL.Path)
	return p
}

func (r *Request) SrchAttr(attr ...string) *Request {
	r.Params.SrchAttr = attr
	return r
}

func (r *Request) FacetAttr(attr ...string) *Request {
	r.Params.FacetAttr = attr
	return r
}

func (r *Request) SortAttr(attr ...string) *Request {
	r.Params.SortAttr = attr
	return r
}

func (r *Request) Facets(attr ...string) *Request {
	r.Params.Facets = attr
	return r
}

func (r *Request) AddFilterVal(field string, filters ...string) *Request {
	for _, f := range filters {
		f = field + ":" + f
		r.filters = append(r.filters, f)
	}
	return r
}

func (r *Request) OrFilter(filters ...string) *Request {
	return r.AddFilters(lo.ToAnySlice(filters))
}

func (r *Request) AndFilter(filters ...string) *Request {
	return r.AddFilters(lo.ToAnySlice(filters)...)
}

func (r *Request) AddFilters(filters ...any) *Request {
	r.Params.FacetFilters = append(r.Params.FacetFilters, filters...)
	return r
}

func NewAnyFilter(field string, filters []string) []any {
	return lo.ToAnySlice(NewFilter(field, filters...))
}

func NewFilter(field string, filters ...string) []string {
	f := make([]string, len(filters))
	for i, filter := range filters {
		f[i] = field + ":" + filter
	}
	return f
}

func (r *Request) FacetFilters(attr ...any) *Request {
	r.Params.FacetFilters = attr
	return r
}

func (r *Request) Query(val string) *Request {
	r.Params.Query = val
	return r
}

func (r *Request) Filters(val string) *Request {
	r.Params.Filters = val
	return r
}

func (r *Request) SortFacetsBy(val string) *Request {
	r.Params.SortFacetsBy = val
	return r
}

func (r *Request) SortBy(val string) *Request {
	r.Params.SortBy = val
	return r
}

func (r *Request) Order(val string) *Request {
	r.Params.Order = val
	return r
}

func (r *Request) DefaultField(val string) *Request {
	r.Params.DefaultField = val
	return r
}

func (r *Request) UID(val string) *Request {
	r.Params.UID = val
	return r
}

func (r *Request) Page(p int) *Request {
	r.Params.Page = p
	return r
}

func (r *Request) HitsPerPage(p int) *Request {
	r.Params.HitsPerPage = p
	return r
}
