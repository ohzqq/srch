package srch

import (
	"fmt"

	"github.com/ohzqq/facet"
)

type Results struct {
	Data   []any          `json:"data"`
	Facets []*facet.Facet `json:"facets"`
	Query  string         `json:"query"`
}

type Result interface {
	fmt.Stringer
}
