package srch

import (
	"encoding/json"
	"net/url"
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

	fils := parseFilters(filters)

	var aor []*roaring.Bitmap
	var bor []*roaring.Bitmap
	for _, filter := range fils {
		field, ok := fields[filter.Facet]
		if !ok {
			break
		}
		val := field.Filter(filter.Value)
		switch filter.Not {
		case true:
			bits = roaring.AndNot(bits, val)
		default:
			switch filter.Operator {
			case And:
				aor = append(aor, val)
			case Or:
				bor = append(bor, val)
			}
		}
	}

	arb := roaring.ParAnd(viper.GetInt("workers"), aor...)
	bits.And(arb)

	orb := roaring.ParOr(viper.GetInt("workers"), bor...)
	bits.Or(orb)

	return bits, nil
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

//func decodeFilter(bits *roaring.Bitmap, fields map[string]*Field, query string) (*roaring.Bitmap, error) {
//  filters, err := unmarshalFilter(query)
//  if err != nil {
//    return nil, err
//  }

//  fils := parseFilters(filters)

//  var aor []*roaring.Bitmap
//  var bor []*roaring.Bitmap
//  for _, filter := range fils {
//    field, ok := fields[filter.Facet]
//    if !ok {
//      break
//    }
//    val := field.Filter(filter.Value)
//    switch filter.Not {
//    case true:
//      bits = roaring.AndNot(bits, val)
//    default:
//      switch filter.Operator {
//      case And:
//        aor = append(aor, val)
//      case Or:
//        bor = append(bor, val)
//      }
//    }
//  }

//  arb := roaring.ParAnd(viper.GetInt("workers"), aor...)
//  bits.And(arb)

//  orb := roaring.ParOr(viper.GetInt("workers"), bor...)
//  bits.Or(orb)

//  return bits, nil
//}

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

func cutFilter(filter string) (string, string, bool) {
	return strings.Cut(filter, ":")
}

func CutFilter(filter string) (string, string, bool) {
	facet, val, _ := strings.Cut(filter, ":")

	if strings.HasPrefix(val, "-") {
		return facet, strings.TrimPrefix(val, "-"), true
	}

	return facet, val, false
}

func parseFilters(filters []any) []filter {
	var and []filter
	for _, fs := range filters {
		switch vals := fs.(type) {
		case string:
			f := filter{Operator: And}
			f.Facet, f.Value, f.Not = CutFilter(vals)
			and = append(and, f)
		case []any:
			for _, val := range cast.ToStringSlice(vals) {
				f := filter{Operator: Or}
				f.Facet, f.Value, f.Not = CutFilter(val)
				and = append(and, f)
			}
		}
	}
	return and
}

func UnmarshalFilterString(filters string) ([]any, error) {
	dec, err := url.QueryUnescape(filters)
	if err != nil {
		return nil, err
	}

	var filter []any
	err = json.Unmarshal([]byte(dec), &filter)
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
