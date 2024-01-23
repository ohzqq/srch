package srch

import (
	"encoding/json"
)

type Response struct {
	*Index
}

func NewResponse(idx *Index) *Response {
	return &Response{
		Index: idx,
	}
}

func (r *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.StringMap())
}

func (r *Response) NbHits() int {
	return int(r.res.GetCardinality())
}

func (r *Response) StringMap() map[string]any {
	idx := New(r.Query.Params)
	idx.Index(r.GetResults())
	m := idx.StringMap()
	m[NbHits] = r.NbHits()
	m[Hits] = idx.Data
	return m
}
