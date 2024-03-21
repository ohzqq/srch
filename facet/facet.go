package facet

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Fields struct {
	params *param.Params
	Facets []*Facet         `json:"facetFields"`
	data   []map[string]any `json:"-"`
	ids    []string
	bits   *roaring.Bitmap
}

func New(data []map[string]any, param *param.Params) (*Fields, error) {
	facets := NewFields(param.Facets)
	facets.params = param
	facets.data = data
	facets.Calculate()

	if len(facets.params.FacetFilters) > 0 {
		filters := facets.params.FacetFilters
		facets.params.FacetFilters = []any{}
		filtered, err := facets.Filter(filters)
		if err != nil {
			return nil, err
		}
		return filtered.Calculate(), nil
	}

	return facets, nil
}

func NewFields(fields []string) *Fields {
	return &Fields{
		bits:   roaring.New(),
		Facets: NewFacets(fields),
	}
}

func (f *Fields) Calculate() *Fields {
	var uid string
	if f.params.UID != "" {
		uid = f.params.UID
	}

	for id, d := range f.data {
		if i, ok := d[uid]; ok {
			id = cast.ToInt(i)
		}
		f.bits.AddInt(id)
		for _, facet := range f.Facets {
			if val, ok := d[facet.attr]; ok {
				facet.Add(val, []int{id})
			}
		}
	}

	for _, facet := range f.Facets {
		facet.Items = facet.Keywords()
		facet.Count = facet.Len()
		facet.Attribute = joinAttr(facet)
	}
	return f
}

func (f *Fields) Filter(filters []any) (*Fields, error) {
	filtered, err := Filter(f.bits, f.Facets, filters)
	if err != nil {
		return nil, err
	}

	f.bits.And(filtered)

	var data []map[string]any
	if f.bits.GetCardinality() > 0 {
		data = f.getHits()
	}

	facets, err := New(data, f.params)
	if err != nil {
		return nil, err
	}

	return facets, nil
}

func (f Fields) getHits() []map[string]any {
	var uid string
	if f.params.UID != "" {
		uid = f.params.UID
	}
	var hits []map[string]any
	for idx, d := range f.data {
		if i, ok := d[uid]; ok {
			idx = cast.ToInt(i)
		}
		if f.bits.ContainsInt(idx) {
			hits = append(hits, d)
		}
	}
	return hits
}

func (f Fields) GetFacet(attr string) *Facet {
	for _, facet := range f.Facets {
		if facet.attr == attr {
			return facet
		}
	}
	return &Facet{}
}

func (f Fields) Len() int {
	return int(f.bits.GetCardinality())
}

func (f *Fields) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Fields) Items() []string {
	var ids []string

	f.bits.Iterate(func(x uint32) bool {
		ids = append(ids, cast.ToString(x))
		return true
	})

	return ids
}

func (f *Fields) MarshalJSON() ([]byte, error) {
	m := make(map[string]int)
	for _, fi := range f.Facets {
		m[fi.attr] = f.Len()
	}
	return json.Marshal(m)
}

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}
