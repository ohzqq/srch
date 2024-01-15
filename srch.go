package srch

import (
	"net/url"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func Filter(data []map[string]any, facets []*Field, values url.Values) []map[string]any {
	var bits []*roaring.Bitmap
	for name, filters := range values {
		for _, facet := range facets {
			if facet.Attribute == name {
				bits = append(bits, facet.Filter(filters...))
			}
		}
	}

	filtered := roaring.ParOr(viper.GetInt("workers"), bits...)
	ids := filtered.ToArray()

	return FilteredItems(data, lo.ToAnySlice(ids))
}

// FilteredItems returns the subset of data.
func FilteredItems(data []map[string]any, ids []any) []map[string]any {
	items := make([]map[string]any, len(ids))
	for item, _ := range data {
		for i, id := range ids {
			if cast.ToInt(id) == item {
				items[i] = data[item]
			}
		}
	}
	return items
}

func GetFieldsFromSlice(items []map[string]any, names []string) []*Field {
	if len(items) < 1 {
		return []*Field{}
	}

	item := items[0]

	if len(names) < 1 {
		names = lo.Keys(item)
	}

	var fields []*Field
	for _, f := range names {
		if _, ok := item[f]; ok {
			fields = append(fields, NewTextField(f))
		}
	}
	return fields
}

func GetSearchableFieldValues(data []map[string]any, fields []string) []string {
	src := make([]string, len(data))
	for i, d := range data {
		s := lo.PickByKeys(d, fields)
		vals := cast.ToStringSlice(lo.Values(s))
		src[i] = strings.Join(vals, "\n")
	}
	return src
}
