package srch

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	TBool = 1
)

func Filter(bits *roaring.Bitmap, fields map[string]*Field, filters []any) (*roaring.Bitmap, error) {
	var and []*roaring.Bitmap
	var or []*roaring.Bitmap
	var not []*roaring.Bitmap
	for name, field := range fields {
		for _, fs := range filters {
			switch vals := fs.(type) {
			case string:
				vals, ok := strings.CutPrefix(vals, name+":")
				if ok {
					vals, n := strings.CutPrefix(vals, "-")
					f := field.Filter(vals)
					if n {
						not = append(not, f)
					} else {
						and = append(and, f)
					}
				}
			case []any:
				os := cast.ToStringSlice(vals)
				for _, o := range os {
					o, ok := strings.CutPrefix(o, name+":")
					if ok {
						o, n := strings.CutPrefix(o, "-")
						f := field.Filter(o)
						if n {
							not = append(not, f)
						} else {
							or = append(or, f)
						}
					}
				}
			}
		}
	}

	for _, n := range not {
		bits.AndNot(n)
	}

	arb := roaring.ParAnd(viper.GetInt("workers"), and...)
	bits.And(arb)

	orb := roaring.ParOr(viper.GetInt("workers"), or...)
	bits.Or(orb)

	return bits, nil
}

func NewAnyFilter(field string, filters []string) []any {
	return lo.ToAnySlice(NewFilter(field, filters...))
}

func NewFilter(field string, filters ...string) []string {
	f := make([]string, len(filters))
	for i, filter := range filters {
		f[i] = field + ":" + filter
	}
	return f
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func bitsToIntSlice(bitmap *roaring.Bitmap) []int {
	bits := bitmap.ToArray()
	ids := make([]int, len(bits))
	for i, b := range bits {
		ids[i] = int(b)
	}
	return ids
}

// FilteredItems returns the subset of data.
func FilteredItems(data []map[string]any, ids []any) []map[string]any {
	if len(ids) == 0 {
		return data
	}
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
