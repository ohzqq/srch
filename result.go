package srch

import (
	"fmt"
	"net/url"
)

type Results struct {
	Data    []any      `json:"data"`
	Facets  []*Facet   `json:"facets"`
	Query   string     `json:"query"`
	Filters url.Values `json:"filters"`
}

type Item interface {
	fmt.Stringer
}
