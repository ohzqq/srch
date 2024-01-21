package srch

type Response struct {
	Hits        []map[string]any `json:"hits"`
	Page        int              `json:"page"`
	NbHits      int              `json:"nbHits"`
	NbPages     int              `json:"nbPages"`
	HitsPerPage int              `json:"hitsPerPage"`
	Keywords    string           `json:"query"`
	Facets      []*Field         `json:"facets"`
	*Query      `json:"params"`
}
