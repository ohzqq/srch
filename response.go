package srch

import (
	"encoding/json"
)

type Response struct {
	*Index
}

func NewResponse(idx *Index) *Response {
	i, err := New(idx.Params.Values)
	if err != nil {
		i = idx
	}
	i.Index(idx.GetResults())

	return &Response{
		Index: i,
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
