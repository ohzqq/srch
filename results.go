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
	Items []*FacetItem
}

// FacetItem is a data structure for a Facet's item.
type FacetItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

func NewFacet(field *Field) *Facet {
	return &Facet{
		Field: field,
		Items: FieldItemsToFacetItems(field.Items),
	}
}

func FieldItemsToFacetItems(fi map[string]*roaring.Bitmap) []*FacetItem {
	items := make([]*FacetItem, len(fi))
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
