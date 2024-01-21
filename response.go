package srch

type Response struct {
	*Index
	Page        int    `json:"page"`
	NbHits      int    `json:"nbHits"`
	NbPages     int    `json:"nbPages"`
	HitsPerPage int    `json:"hitsPerPage"`
	Keywords    string `json:"query"`
}
