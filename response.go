package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"slices"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Response struct {
	*Index
	facets map[string]*Field
}

func NewResponse(data []map[string]any, vals url.Values) *Response {
	idx, err := New(vals)
	if err != nil {
		log.Fatal(err)
	}
	idx.Data = data
	idx.res = idx.Bitmap()
	//idx.Index(data)
	r := &Response{
		Index:  idx,
		facets: idx.Params.Facets(),
	}

	r.calculateFacets()

	return r
}

func (idx *Response) calculateFacets() {
	for id, d := range idx.Data {
		for _, attr := range idx.FacetAttr() {
			if val, ok := d[attr]; ok {
				idx.facets[attr].Add(val, []int{id})
			}
		}
	}
}

func (idx *Response) Filter(q string) *Response {
	//println(q)
	//filters := idx.Filters()
	filters, err := unmarshalFilter(q)
	if err != nil {
		return idx
	}

	filtered, err := Filter(idx.res, idx.facets, filters)
	if err != nil {
		fmt.Printf("%+v\n", filters)
		return idx
	}

	idx.res.And(filtered)
	return idx.Response()
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
		ParamFacets:        r.Facets(),
	}

	page := r.Page()
	hpp := r.HitsPerPage()
	nbh := r.NbHits()
	m[HitsPerPage] = hpp
	m[NbHits] = nbh
	m[Page] = page

	if nbh > 0 {
		m["nbPages"] = nbh/hpp + 1
	}

	m[Hits] = r.VisibleHits(page, nbh, hpp)

	return m
}

func (r *Response) VisibleHits(page, nbh, hpp int) []map[string]any {
	if nbh < hpp {
		return r.Data
	}
	b := hpp * page
	e := b + hpp
	return lo.Slice(r.Data, b, e)
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

func sortDataByTextField(data []map[string]any, field string) []map[string]any {
	fn := func(a, b map[string]any) int {
		x := cast.ToString(a[field])
		y := cast.ToString(b[field])
		switch {
		case x > y:
			return 1
		case x == y:
			return 0
		default:
			return -1
		}
	}
	slices.SortFunc(data, fn)
	return data
}

func sortDataByNumField(data []map[string]any, field string) []map[string]any {
	fn := func(a, b map[string]any) int {
		x := cast.ToInt(a[field])
		y := cast.ToInt(b[field])
		switch {
		case x > y:
			return 1
		case x == y:
			return 0
		default:
			return -1
		}
	}
	slices.SortFunc(data, fn)
	return data
}
