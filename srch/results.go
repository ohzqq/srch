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

func NewResults(hits []map[string]any, params *param.Params) *Results {
	return &Results{
		hits:   hits,
		Params: params,
	}
}
