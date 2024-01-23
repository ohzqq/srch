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
	m := r.StringMap()
	d, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (r *Response) NbHits() int {
	return int(r.res.GetCardinality())
}

func (r *Response) StringMap() map[string]any {
	m := make(map[string]any)
	idx := New(r.Query.Params)
	idx.Index(r.GetResults())
	m[NbHits] = r.NbHits()
	m[ParamQuery] = r.Query.Query()
	m[Page] = idx.Page()
	m[Hits] = idx.Data
	m["params"] = r.Query
	m[HitsPerPage] = idx.HitsPerPage()
	//if idx.Query.Has(ParamFacets)
	m[ParamFacets] = idx.Facets()
	return m
}
