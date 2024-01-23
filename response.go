package srch

import (
	"encoding/json"
)

type Response struct {
	*Index
	Page        int    `json:"page"`
	NbPages     int    `json:"nbPages"`
	HitsPerPage int    `json:"hitsPerPage"`
	Keywords    string `json:"query"`
}

func NewResponse(idx *Index) *Response {
	return &Response{
		Index: idx,
	}
}

func (r *Response) MarshalJSON() ([]byte, error) {
	m := IndexToResponseMap(r.Index)
	d, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (r *Response) NbHits() int {
	return int(r.res.GetCardinality())
}

func IndexToResponseMap(idx *Index) map[string]any {
	m := make(map[string]any)
	m[NbHits] = idx.Len()
	m[ParamQuery] = idx.Query.Query()
	m[Page] = idx.Page()
	m[Hits] = idx.getDataByBitmap(idx.res)
	m["params"] = idx.Query.Params
	m[HitsPerPage] = idx.HitsPerPage()
	m[ParamFacets] = idx.Facets()
	return m
}
