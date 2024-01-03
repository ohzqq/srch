package srch

import (
	"net/url"
)

type Results struct {
	*Src
	Facets  []*Facet   `json:"facets"`
	Filters url.Values `json:"filters"`
}

func NewResults(data []any) *Results {
	return &Results{
		Src: NewSrc(data),
	}
}
