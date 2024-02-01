package srch

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
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
	m := map[string]any{
		"processingTimeMS": 1,
		"params":           r.Params,
		Query:              r.Params.Query(),
		Page:               r.Page(),
		ParamFacets:        r.Facets(),
		Hits:               r.Data,
	}

	hpp := r.HitsPerPage()
	nbh := r.NbHits()
	m[HitsPerPage] = hpp
	m[NbHits] = nbh

	if nbh > 0 {
		m["nbPages"] = nbh/hpp + 1
	}
	return m
}

// JSON marshals an Index to json.
func (idx *Response) JSON() []byte {
	d, err := json.Marshal(idx)
	if err != nil {
		return []byte{}
	}
	return d
}

// Print writes Index json to stdout.
func (idx *Response) Print() {
	enc := json.NewEncoder(os.Stdout)
	err := enc.Encode(idx)
	if err != nil {
		log.Fatal(err)
	}
}

// PrettyPrint writes Index indented json to stdout.
func (idx *Response) PrettyPrint() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(idx)
	if err != nil {
		log.Fatal(err)
	}
}
