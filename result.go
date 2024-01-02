package srch

import (
	"fmt"
)

type Results struct {
	Data   []any    `json:"data"`
	Facets []*Facet `json:"facets"`
	Query  string   `json:"query"`
}

type Item interface {
	fmt.Stringer
}
