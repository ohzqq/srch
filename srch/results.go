package srch

import (
	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
)

type Results struct {
	Facets *facet.Facets
	Params *param.Params
	hits   []map[string]any
}

func NewResults(hits []map[string]any, params *param.Params) (*Results, error) {
	res := &Results{
		hits:   hits,
		Params: params,
	}

	if len(hits) == 0 {
		return res, nil
	}

	if len(params.Facets) > 0 {
		facets, err := facet.New(hits, params.FacetSettings)
		if err != nil {
			return nil, err
		}
		res.Facets = facets
	}

	return res, nil
}

func (res *Results) NbHits() int {
	return len(res.hits)
}

func (r *Results) StringMap() map[string]any {

	m := map[string]any{
		"processingTimeMS": 1,
		"params":           r.Params,
		param.Query:        r.Params.Query,
		param.Facets:       r.Params.Facets,
	}

	page := r.Page()
	hpp := r.HitsPerPage()
	nbh := r.NbHits()
	m[HitsPerPage] = hpp
	m[NbHits] = nbh
	m[Page] = page

	if nbh > 0 {
		m["nbPages"] = nbh/hpp + 1
	}

	m[Hits] = r.VisibleHits(page, nbh, hpp)

	return m
}
