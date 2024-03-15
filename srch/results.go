package srch

import (
	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
)

type Results struct {
	Facets *facet.Facets
	Params *param.Params
}
