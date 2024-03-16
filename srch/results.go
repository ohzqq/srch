package srch

import (
	"fmt"

	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
)

type Results struct {
	Facets *facet.Facets
	Params *param.Params
	nbHits []map[string]any
}

func NewResults(hits []map[string]any, params *param.Params) (*Results, error) {
	res := &Results{
		nbHits: hits,
		Params: params,
	}

	if len(hits) == 0 {
		return res, nil
	}

	if params.Has(param.Facets) {
		facets, err := facet.New(hits, params.FacetSettings)
		if err != nil {
			return nil, err
		}
		res.Facets = facets
	}

	return res, nil
}

func (res *Results) NbHits() int {
	return len(res.nbHits)
}

func (res *Results) Hits() []map[string]any {
	if !res.Params.Has(param.HitsPerPage) {
		return res.nbHits
	}
	nbHits := res.NbHits()
	hpp := res.Params.HitsPerPage
	if nbHits < hpp {
		return res.nbHits
	}
	fmt.Printf("hpp %d\n", hpp)
	nbPages := nbHits/hpp + nbHits%hpp
	fmt.Printf("nbPages %d\n", nbPages)
	page := 0
	if res.Params.Has(param.Page) {
		page = res.Params.Page - 1
	}
	fmt.Printf("page %d\n", page)

	return res.nbHits
}

//func (r *Results) StringMap() map[string]any {
//  m := map[string]any{
//    "processingTimeMS": 1,
//    "params":           r.Params,
//    param.Query:        r.Params.Query,
//    param.Facets:       r.Params.Facets,
//  }

//  page := r.Page()
//  hpp := r.HitsPerPage()
//  nbh := r.NbHits()
//  m[HitsPerPage] = hpp
//  m[NbHits] = nbh
//  m[Page] = page

//  if nbh > 0 {
//    m["nbPages"] = nbh/hpp + 1
//  }

//  m[Hits] = r.VisibleHits(page, nbh, hpp)

//  return m
//}
