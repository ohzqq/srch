package srch

import "net/url"

type Response struct {
	*Index
	Page        int    `json:"page"`
	NbHits      int    `json:"nbHits"`
	NbPages     int    `json:"nbPages"`
	HitsPerPage int    `json:"hitsPerPage"`
	Keywords    string `json:"query"`
}

func NewResponse(data []map[string]any, params url.Values) *Response {
	res := &Response{
		Index:    New(params).Index(data),
		Page:     0,
		Keywords: params.Get(ParamQuery),
	}
	res.NbHits = res.Len()
	return res
}
