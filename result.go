package srch

import "net/url"

type Results struct {
	Data    []any      `json:"data"`
	Facets  []*Facet   `json:"facets"`
	Filters url.Values `json:"filters"`
}
