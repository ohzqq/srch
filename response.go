package srch

import (
	"encoding/json"
	"log"
	"net/url"
)

type Response struct {
	*Index
}

func NewResponse(data []map[string]any, vals url.Values) *Response {
	idx, err := New(vals)
	if err != nil {
		log.Fatal(err)
	}
	idx.Index(data)
	return &Response{
		Index: idx,
	}
}

func (r *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.StringMap())
}

func (r *Response) NbHits() int {
	if r.HasResults() {
		return int(r.res.GetCardinality())
	}
	return int(r.Index.Bitmap().GetCardinality())
}

func (r *Response) StringMap() map[string]any {
	m := r.StringMap()
	m[NbHits] = r.NbHits()
	m[Hits] = r.Data
	return m
}
