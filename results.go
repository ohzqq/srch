package srch

import (
	"github.com/RoaringBitmap/roaring"
)

type Results struct {
	Data   []map[string]any
	Facets []*Facet
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

func NewResults(data []map[string]any) *Results {
	return &Results{
		Data: data,
	}
}

func (r *Results) Filter(q string) *Results {
	vals, err := ParseValues(q)
	if err != nil {
		return r
	}
	r.SetData(Filter(r.Data, FacetsToFields(r.Facets), vals))
	return r
}

func (r *Results) SetFacets(facets []*Field) *Results {
	r.Facets = FieldsToFacets(facets)
	return r
}

func (r *Results) SetData(data []map[string]any) *Results {
	r.Data = data
	return r
}

func (r *Results) Src() []map[string]any {
	return r.Data
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
