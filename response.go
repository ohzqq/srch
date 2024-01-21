package srch

import "net/url"

type Response struct {
	Hits        []map[string]any
	Page        int
	NbHits      int
	NbPages     int
	HitsPerPage int
	Query       string
	Params      url.Values
	Facets      map[string]int
}
