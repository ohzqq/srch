package srch

import (
	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Response struct {
	Facets []*facet.Field
	Params *param.Params
	hits   []map[string]any
}

func NewResponse(hits []map[string]any, params *param.Params) (*Response, error) {
	res := &Response{
		hits:   hits,
		Params: params,
	}

	if len(hits) == 0 {
		return res, nil
	}

	if params.Has(param.Facets) {
		//h := FilterDataByAttr(hits, params.Facets)
		facets, err := facet.New(hits, params)
		if err != nil {
			return nil, err
		}
		res.Facets = facets.Fields
	}

	return res, nil
}

func (res *Response) NbHits() int {
	return len(res.hits)
}

func (res *Response) NbPages() int {
	hpp := 1
	if res.Params.Has(param.HitsPerPage) {
		hpp = res.Params.HitsPerPage
	}

	nb := res.NbHits() / hpp
	if r := res.NbHits() % hpp; r > 0 {
		nb++
	}

	return nb
}

func (res *Response) HitsPerPage() int {
	if res.Params.Has(param.HitsPerPage) {
		return res.Params.HitsPerPage
	}
	return viper.GetInt("hitsPerPage")
}

func (res *Response) Page() int {
	if !res.Params.Has(param.Page) {
		return 0
	}
	return res.Params.Page
}

func (res *Response) page() int {
	return res.Page() - 1
}

func (res *Response) Hits() []map[string]any {
	nbHits := res.NbHits()
	hpp := res.HitsPerPage()

	if nbHits < hpp {
		return res.hits
	}

	page := res.Page()

	if nb := res.NbPages(); page >= nb {
		return []map[string]any{}
	}

	return lo.Subset(res.hits, page*hpp, uint(hpp))
}

func (r *Response) StringMap() map[string]any {
	m := map[string]any{
		"processingTimeMS": 1,
		"params":           r.Params.Encode(),
		param.Query:        r.Params.Query,
		param.Facets:       r.Facets,
	}

	page := r.Page()
	hpp := r.HitsPerPage()
	nbh := r.NbHits()
	m[param.HitsPerPage] = hpp
	m[param.NbHits] = nbh
	m[param.Page] = page
	m[param.NbPages] = r.NbPages()

	m[param.Hits] = r.Hits()

	return m
}
