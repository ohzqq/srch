package srch

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type filter struct {
	Facet    string
	Operator string
	Not      bool
	Value    string
}

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
			//for _, v := range vals {
			//  not, ok := IsNegative(v)
			//  if ok {
			//    bits.AndNot(facet.Filter(not))
			//  } else {
			//  }
			//}

		}
	}
	arb := roaring.ParAnd(viper.GetInt("workers"), aor...)
	bits.And(arb)

	orb := roaring.ParOr(viper.GetInt("workers"), bor...)
	bits.Or(orb)

	return bits, nil
}

func parseFilters(filters []any) ([]string, []string) {
	var and, or []string
	for _, fs := range filters {
		switch vals := fs.(type) {
		case string:
			and = append(and, vals)
		case []any:
			or = append(or, cast.ToStringSlice(vals)...)
		}
	}
	return and, or
}

func CutFilter(filter string) (string, string, bool) {
	facet, val, _ := strings.Cut(filter, ":")

	if strings.HasPrefix(val, "-") {
		return facet, strings.TrimPrefix(val, "-"), true
	}

	return facet, val, false
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

func IsNegative(filter string) (string, bool) {
	return strings.TrimPrefix(filter, "-"), strings.HasPrefix(filter, "-")
}

func FilterByAttribute(attr string, filters []string) []string {
	fn := func(f string, _ int) (string, bool) {
		pre := attr + ":"
		return strings.TrimPrefix(f, pre), strings.HasPrefix(f, pre)
	}
	return lo.FilterMap(filters, fn)
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
