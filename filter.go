package srch

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Filters struct {
	Con     url.Values
	Dis     url.Values
	filters map[string]url.Values
	labels  map[string]url.Values
}

func Filter(bits *roaring.Bitmap, fields []*Field, filters *Filters) (*roaring.Bitmap, error) {
	for _, facet := range fields {
		filters.Filter(bits, facet)
	}

	return bits, nil
}

func (f *Filters) Filter(bits *roaring.Bitmap, facet *Field) {
	if filter, ok := f.labels[facet.Attribute]; ok {
		for op, vals := range filter {
			for _, v := range vals {
				not, ok := IsNegative(v)
				if ok {
					bits.AndNot(facet.Filter(not))
				} else {
					switch op {
					case And:
						bits.And(facet.Filter(v))
					case Or:
						bits.Or(facet.Filter(v))
					}
				}
			}
		}
	}

	//if f.Con.Has(facet.Attribute) {
	//  for _, a := range f.Con[facet.Attribute] {
	//    not, ok := IsNegative(a)
	//    if ok {
	//      bits.AndNot(facet.Filter(not))
	//    } else {
	//      bits.And(facet.Filter(a))
	//    }
	//  }
	//}

	//if f.Dis.Has(facet.Attribute) {
	//  for _, or := range f.Dis[facet.Attribute] {
	//    not, ok := IsNegative(or)
	//    if ok {
	//      bits.AndNot(facet.Filter(not))
	//    } else {
	//      bits.Or(facet.Filter(or))
	//    }
	//  }
	//}

}

func DecodeFilter(query string, facets ...string) (*Filters, error) {
	filters := &Filters{
		Con:    make(url.Values),
		Dis:    make(url.Values),
		labels: make(map[string]url.Values),
	}

	for _, facet := range facets {
		filters.labels[facet] = make(url.Values)
	}

	ff, err := unmarshalFilter(query)
	if err != nil {
		return nil, err
	}

	for _, v := range ff {
		switch vals := v.(type) {
		case string:
			filters.addCon(vals)
		case []any:
			or := cast.ToStringSlice(vals)
			switch len(or) {
			case 1:
				filters.addCon(or[0])
			default:
				filters.addDis(or)
			}
		}
	}

	return filters, nil
}

func (f *Filters) addCon(fv string) {
	label, filter, ok := cutFilter(fv)
	if !ok {
		return
	}
	f.labels[label].Add(And, filter)
	f.Con.Add(label, filter)
}

func (f *Filters) addDis(filters []string) {
	for _, fv := range filters {
		label, filter, ok := cutFilter(fv)
		if !ok {
			break
		}
		f.labels[label].Add(Or, filter)
		f.Dis.Add(label, filter)
	}
}

func (f *Filters) Encode() string {
	return f.ToValues().Encode()
}

func (f *Filters) String() string {
	return string(f.Bytes())
}

func (f *Filters) Bytes() []byte {
	var filters []any
	for k, and := range f.Con {
		filters = append(filters, mapFilterVals(k, and)...)
	}
	for k, or := range f.Dis {
		filters = append(filters, mapFilterVals(k, or))
	}

	filter, err := json.Marshal(filters)
	if err != nil {
		filter = []byte{}
	}

	return filter
}

func (f *Filters) ToValues() url.Values {
	return url.Values{
		"facetFilters": []string{f.String()},
	}
}

func mapFilterVals(key string, vals []string) []any {
	m := make([]any, len(vals))
	for i, v := range vals {
		m[i] = key + ":" + v
	}
	return m
}

func IsNegative(f string) (string, bool) {
	return strings.TrimPrefix(f, "-"), strings.HasPrefix(f, "-")
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

func UnmarshalFilterString(filters string) ([]any, error) {
	dec, err := url.QueryUnescape(filters)
	if err != nil {
		return nil, err
	}

	return unmarshalFilter(dec)
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func ParseFilters(filters any) (string, []string, error) {
	switch vals := filters.(type) {
	case string:
		return And, []string{vals}, nil
	case []any:
		return Or, cast.ToStringSlice(vals), nil
	default:
		return "", []string{}, errors.New("not a filter")
	}
}
