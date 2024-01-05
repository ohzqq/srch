package srch

import (
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Filter takes an *Index, filters the data and calculates the facets. It
// returns a new *Index.
func Filter(idx *Index) *Index {
	for _, f := range idx.Facets {
		fmt.Printf("%+v\n", f)
	}
	var bits []*roaring.Bitmap
	for name, filters := range idx.Query {
		for _, facet := range idx.Facets {
			if facet.Attribute == name {
				bits = append(bits, facet.Filter(filters...))
			}
		}
	}

	filtered := roaring.ParOr(viper.GetInt("workers"), bits...)
	ids := filtered.ToArray()

	d := FilteredItems(idx.Data, lo.ToAnySlice(ids))
	println(len(d))
	return CopyIndex(idx, d)
}

// FilteredItems returns the subset of data.
func FilteredItems(data []any, ids []any) []any {
	items := make([]any, len(ids))
	for item, _ := range data {
		for i, id := range ids {
			if cast.ToInt(id) == item {
				items[i] = data[item]
			}
		}
	}
	return items
}
