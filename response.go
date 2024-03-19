package srch

import (
	"fmt"

	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Response struct {
	*param.Params

	results []map[string]any

	RawQuery string           `json:"params"`
	Facets   []*facet.Field   `json:"facetFields"`
	FacetMap map[string]int   `json:"facets"`
	Hits     []map[string]any `json:"hits"`
	NbHits   int              `json:"nbHits"`
	NbPages  int              `json:"nbPages"`
}

func NewResponse(hits []map[string]any, params *param.Params) (*Response, error) {
	res := &Response{
		results:  hits,
		Params:   params,
		RawQuery: params.Encode(),
		FacetMap: make(map[string]int),
	}

	if len(hits) == 0 {
		return res, nil
	}

	res.NbHits = res.nbHits()

	if params.Has(param.Facets) {
		facets, err := facet.New(hits, params)
		if err != nil {
			return nil, fmt.Errorf("response failed to calculate facets: %w\n", err)
		}
		res.Facets = facets.Fields
		for _, f := range facets.Fields {
			res.FacetMap[f.Attribute] = f.Len()
		}
	}

	res.calculatePagination()

	res.Hits = res.visibleHits()

	return res, nil
}

func (res *Response) calculatePagination() *Response {
	res.HitsPerPage = res.hitsPerPage()
	res.Page = res.page()
	res.NbPages = res.nbPages()
	return res
}

func (res *Response) nbHits() int {
	return len(res.results)
}

func (res *Response) nbPages() int {
	hpp := res.HitsPerPage

	nb := res.NbHits / hpp
	if r := res.NbHits % hpp; r > 0 {
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

func (res *Response) visibleHits() []map[string]any {
	if res.NbHits < res.HitsPerPage {
		return res.results
	}

	if nb := res.NbPages; res.Page >= nb {
		return []map[string]any{}
	}

	return lo.Subset(res.results, res.Page*res.HitsPerPage, uint(res.HitsPerPage))
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

	return m
}
