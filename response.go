package srch

import "net/url"

type Response struct {
	Hits        []map[string]any `json:"hits"`
	Page        int              `json:"page"`
	NbHits      int              `json:"nbHits"`
	NbPages     int              `json:"nbPages"`
	HitsPerPage int              `json:"hitsPerPage"`
	Query       string           `json:"query"`
	Params      url.Values       `json:"params"`
	Facets      []*Field         `json:"facets"`
}
