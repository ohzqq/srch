package srch

import (
	"fmt"

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
	println(len(hits))

	fmt.Printf("facets %#v\n", params.Facets)
	fmt.Printf("hits %#v\n", hits[0])
	//fmt.Printf("facets %#v\n", maps.Keys(hits[0]))

	if len(params.Facets) > 0 {
		facets, err := facet.New(hits, params.FacetSettings)
		if err != nil {
			return nil, err
		}
		res.Facets = facets
	}

	return res, nil
}
