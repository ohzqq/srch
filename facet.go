package srch

import (
	"encoding/json"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
)

type Facet struct {
	*Field
	Items []*FacetItem `json:"items"`
}

// FacetItem is a data structure for a Facet's item.
type FacetItem struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	bits        *roaring.Bitmap
	fuzzy.Match `json:"-"`
}

func (f *FacetItem) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		"value": f.Value,
		"label": f.Label,
		"count": f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (f *FacetItem) Count() int {
	return len(f.bits.ToArray())
}

func NewFacet(field *Field) *Facet {
	f := &Facet{
		Field: field,
		Items: field.ItemsWithCount(),
	}
	return f
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

func SortItemsByCount(items []*FacetItem) []*FacetItem {
	slices.SortFunc(items, sortByCountFunc)
	return items
}

func SortItemsByLabel(items []*FacetItem) []*FacetItem {
	slices.SortFunc(items, sortByLabelFunc)
	return items
}

func sortByCountFunc(a *FacetItem, b *FacetItem) int {
	aC := a.Count()
	bC := b.Count()
	switch {
	case aC < bC:
		return 1
	case aC == bC:
		return 0
	default:
		return -1
	}
}

func sortByLabelFunc(a *FacetItem, b *FacetItem) int {
	switch {
	case a.Label < b.Label:
		return 1
	case a.Label == b.Label:
		return 0
	default:
		return -1
	}
}

func NewFacetItem(label string) *FacetItem {
	return &FacetItem{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

// GetItem returns an *FacetItem.
func (f *Facet) GetItem(term string) *FacetItem {
	for _, item := range f.Items {
		if term == item.Value {
			return item
		}
	}
	return &FacetItem{}
}

// ListItems returns a string slice of all item values.
func (f *Facet) ListItems() []string {
	var items []string
	for _, item := range f.Items {
		items = append(items, item.Value)
	}
	return items
}

// FuzzyFindItem fuzzy finds an item's value and returns possible matches.
func (f *Facet) FuzzyFindItem(term string) []*FacetItem {
	matches := f.FuzzyMatches(term)
	items := make([]*FacetItem, len(matches))
	for i, match := range matches {
		item := f.Items[match.Index]
		item.Match = match
		items[i] = item
	}
	return items
}

// FuzzyMatches returns the fuzzy.Matches of the search.
func (f *Facet) FuzzyMatches(term string) fuzzy.Matches {
	return fuzzy.FindFrom(term, f)
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (f *Facet) String(i int) string {
	return f.Items[i].Value
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Facet) Len() int {
	return len(f.Items)
}
