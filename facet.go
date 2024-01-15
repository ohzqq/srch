package srch

import (
	"encoding/json"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
)

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
