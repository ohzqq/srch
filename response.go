package srch

import (
	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Response struct {
	*param.Params
	Facets []*facet.Field   `json:"facets"`
	hits   []map[string]any `json:"hits"`
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
		facets, err := facet.New(hits, params)
		if err != nil {
			return nil, err
		}
		res.Facets = facets.Fields
	}

	return res, nil
}

func (res *Response) nbHits() int {
	return len(res.hits)
}

func (res *Response) nbPages() int {
	hpp := res.hitsPerPage()

	nb := res.nbHits() / hpp
	if r := res.nbHits() % hpp; r > 0 {
		nb++
	}

	return nb
}

func (res *Response) hitsPerPage() int {
	if res.Params.Has(param.HitsPerPage) {
		return res.Params.HitsPerPage
	}
	return viper.GetInt("hitsPerPage")
}

func (res *Response) page() int {
	if !res.Params.Has(param.Page) {
		return 0
	}
	return res.Params.Page
}

func (res *Response) Hits() []map[string]any {
	nbHits := res.nbHits()
	hpp := res.hitsPerPage()

	if nbHits < hpp {
		return res.hits
	}

	page := res.page()

	if nb := res.nbPages(); page >= nb {
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

	page := r.page()
	hpp := r.hitsPerPage()
	nbh := r.nbHits()
	m[param.HitsPerPage] = hpp
	m[param.NbHits] = nbh
	m[param.Page] = page
	m[param.NbPages] = r.nbPages()

	m[param.Hits] = r.Hits()

	return m
}
