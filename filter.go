package srch

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type filter struct {
	Facet    string
	Operator string
	Not      bool
	Value    string
}

type Filters []any

func Filter(bits *roaring.Bitmap, fields map[string]*Field, query string) (*roaring.Bitmap, error) {
	filters, err := unmarshalFilter(query)
	if err != nil {
		return nil, err
	}

	var aor []*roaring.Bitmap
	var bor []*roaring.Bitmap
	for name, field := range fields {
		for _, fs := range filters {
			switch vals := fs.(type) {
			case string:
				vals, ok := strings.CutPrefix(vals, name+":")
				if ok {
					vals, not := strings.CutPrefix(vals, "-")
					if not {
						bits.AndNot(field.Filter(vals))
					} else {
						aor = append(aor, field.Filter(vals))
					}
				}
			case []any:
				or := cast.ToStringSlice(vals)
				for _, o := range or {
					o, ok := strings.CutPrefix(o, name+":")
					if ok {
						o, not := strings.CutPrefix(o, "-")
						if not {
							bits.AndNot(field.Filter(o))
						} else {
							bor = append(bor, field.Filter(o))
						}
					}
				}
			}
		}
	}
	arb := roaring.ParAnd(viper.GetInt("workers"), aor...)
	bits.And(arb)

	orb := roaring.ParOr(viper.GetInt("workers"), bor...)
	bits.Or(orb)

	return bits, nil
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
