package facet

import (
	"encoding/json"
	"io"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Facets struct {
	params  *param.FacetSettings
	*Params `json:"params"`
	Facets  []*Field         `json:"facets"`
	data    []map[string]any `json:"hits"`
	ids     []string
	bits    *roaring.Bitmap
}

func New(data []map[string]any, param *param.FacetSettings) (*Facets, error) {
	facets := &Facets{
		params: param,
		Facets: NewFields(param.Facets),
		data:   data,
		bits:   roaring.New(),
	}
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

func NewFacets(fields []string) *Facets {
	return &Facets{
		bits:   roaring.New(),
		Facets: NewFields(fields),
	}
}

func (f *Facets) Calculate() *Facets {
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
			if val, ok := d[facet.Attribute]; ok {
				facet.Add(
					val,
					[]int{id},
				)
			}
		}
	}
	return f
}

func (f *Facets) Filter(filters []any) (*Facets, error) {
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

func (f Facets) getHits() []map[string]any {
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

func (f Facets) GetFacet(attr string) *Field {
	for _, facet := range f.Facets {
		if facet.Attribute == attr {
			return facet
		}
	}
	return &Field{}
}

func (f Facets) Len() int {
	return int(f.bits.GetCardinality())
}

func (f Facets) EncodeQuery() string {
	return f.vals.Encode()
}

func (f *Facets) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Facets) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	for _, d := range f.data {
		err := enc.Encode(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Facets) MarshalJSON() ([]byte, error) {
	enc := f.resultMeta()
	enc["hits"] = f.data
	if len(f.data) < 1 {
		enc["hits"] = []any{}
	}

	return json.Marshal(enc)
}

func (f *Facets) resultMeta() map[string]any {
	enc := make(map[string]any)

	facets := make(map[string]*Field)
	for _, facet := range f.Facets {
		facets[facet.Attribute] = facet
	}
	enc["facets"] = facets

	f.vals.Set("nbHits", cast.ToString(f.Len()))
	enc["params"] = f.EncodeQuery()
	return enc
}

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}
