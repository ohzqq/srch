package srch

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Results struct {
	idx              *Index
	Data             []map[string]any `json:"data"`
	Facets           []*Facet         `json:"facets"`
	searchableFields []string
}

type Facet struct {
	*Field
	Items []*FacetItem `json:"items"`
}

// FacetItem is a data structure for a Facet's item.
type FacetItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

func NewResults(idx *Index, data []map[string]any) *Results {
	res := &Results{
		idx:  idx,
		Data: data,
	}
	if idx.HasFacets() {
		res.Facets = FieldsToFacets(idx.Facets())
	}
	return res
}

func (r *Results) Filter(q string) *Results {
	vals, err := ParseValues(q)
	if err != nil {
		return r
	}
	r.SetData(Filter(r.Data, r.idx.Facets(), vals))
	return r
}

func (r *Results) SetData(data []map[string]any) *Results {
	r.Data = data
	return r
}

func (r *Results) Src() []map[string]any {
	return r.Data
}

func (r *Results) Choose() (*Results, error) {
	ids, err := Choose(r)
	if err != nil {
		return r, err
	}

	r.Data = collectResults(r.Data, ids)

	return r, nil
}

func (r *Results) String(i int) string {
	s := lo.PickByKeys(
		r.Data[i],
		r.idx.SearchableFields(),
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}

func (r *Results) Len() int {
	return len(r.Data)
}

func NewFacet(field *Field) *Facet {
	return &Facet{
		Field: field,
		Items: FieldItemsToFacetItems(field.Items),
	}
}

func FieldsToFacets(fields []*Field) []*Facet {
	facets := make([]*Facet, len(fields))
	for i, f := range fields {
		facets[i] = NewFacet(f)
	}
	return facets
}

func FacetsToFields(fields []*Facet) []*Field {
	facets := make([]*Field, len(fields))
	for i, f := range fields {
		facets[i] = f.Field
	}
	return facets
}

func FieldItemsToFacetItems(fi map[string]*roaring.Bitmap) []*FacetItem {
	var items []*FacetItem
	for label, bits := range fi {
		items = append(items, NewFacetItem(label, len(bits.ToArray())))
	}
	return items
}

// NewFacetItem initializes an item with a value and string slice of related data
// items.
func NewFacetItem(name string, count int) *FacetItem {
	return &FacetItem{
		Value: name,
		Label: name,
		Count: count,
	}
}
