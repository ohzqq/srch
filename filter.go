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
	Con url.Values
	Dis url.Values
	Neg url.Values
}

func Filter(bits *roaring.Bitmap, fields []*Field, filters *Filters) (*roaring.Bitmap, error) {
	for _, facet := range fields {
		if filters.Con.Has(facet.Attribute) {
			for _, a := range filters.Con[facet.Attribute] {
				not, ok := IsNegative(a)
				if ok {
					bits.AndNot(facet.Filter(not))
				} else {
					bits.And(facet.Filter(a))
				}
			}
		}
		if filters.Dis.Has(facet.Attribute) {
			for _, or := range filters.Dis[facet.Attribute] {
				not, ok := IsNegative(or)
				if ok {
					bits.AndNot(facet.Filter(not))
				} else {
					bits.Or(facet.Filter(or))
				}
			}
		}
	}

	return bits, nil
}

func DecodeFilter(query string) (*Filters, error) {
	filters := &Filters{
		Neg: make(url.Values),
		Con: make(url.Values),
		Dis: make(url.Values),
	}

	ff, err := unmarshalFilter(query)
	if err != nil {
		return nil, err
	}

	for _, v := range ff {
		filters.add(v)
	}

	return filters, nil
}

func (f *Filters) add(filters any) {
	switch vals := filters.(type) {
	case string:
		f.And(vals)
	case []any:
		or := cast.ToStringSlice(vals)
		switch len(or) {
		case 1:
			f.And(or[0])
		default:
			for _, filter := range or {
				f.Or(filter)
			}
		}
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
	for k, not := range f.Neg {
		for _, n := range not {
			filters = append(filters, k+":-"+n)
		}
	}
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

func (f *Filters) And(fv string) *Filters {
	label, filter, ok := cutFilter(fv)
	if !ok {
		return f
	}
	//if strings.HasPrefix(filter, "-") {
	//  return f.Not(label, filter)
	//}
	f.Con.Add(label, filter)
	return f
}

func (f *Filters) Or(fv string) *Filters {
	label, filter, ok := cutFilter(fv)
	if !ok {
		return f
	}
	//if strings.HasPrefix(filter, "-") {
	//  return f.Not(label, filter)
	//}
	f.Dis.Add(label, filter)
	return f
}

func IsNegative(f string) (string, bool) {
	return strings.TrimPrefix(f, "-"), strings.HasPrefix(f, "-")
}

func (f *Filters) Not(label, filter string) *Filters {
	filter = strings.TrimPrefix(filter, "-")
	f.Neg.Add(label, filter)
	return f
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
