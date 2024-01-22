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

type Filters struct {
	not []string
	and []string
	or  []string
	Con url.Values
	Dis url.Values
	Neg url.Values
}

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

func newFilters() *Filters {
	return &Filters{
		Neg: make(url.Values),
		Con: make(url.Values),
		Dis: make(url.Values),
	}
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

func DecodeFilter(query string) (*Filters, error) {
	filters, err := UnmarshalFilterString(query)
	if err != nil {
		return nil, err
	}

	f := newFilters()
	for _, v := range filters {
		switch val := v.(type) {
		case string:
			f.add(val)
		case []any:
			f.add(val...)
		}
	}
	return f, nil
}

func cutFilter(filter string) (string, string, bool) {
	return strings.Cut(filter, ":")
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

func (f *Filters) add(vals ...any) *Filters {
	switch filters := cast.ToStringSlice(vals); len(filters) {
	case 1:
		return f.And(filters[0])
	default:
		for _, filter := range filters {
			f.Or(filter)
		}
		return f
	}
}

func (f *Filters) And(fv string) *Filters {
	label, filter, ok := cutFilter(fv)
	if !ok {
		return f
	}
	if strings.HasPrefix(filter, "-") {
		return f.Not(label, filter)
	}
	f.Con.Add(label, filter)
	return f
}

func (f *Filters) Or(fv string) *Filters {
	label, filter, ok := cutFilter(fv)
	if !ok {
		return f
	}
	if strings.HasPrefix(filter, "-") {
		return f.Not(label, filter)
	}
	f.Dis.Add(label, filter)
	return f
}

func (f *Filters) Not(label, filter string) *Filters {
	filter = strings.TrimPrefix(filter, "-")
	f.Neg.Add(label, filter)
	return f
}

func (f *Filters) Encode() string {
	return f.ToValues().Encode()
}

func (f *Filters) String() string {
	return string(f.Bytes())
}

func (f *Filters) Bytes() []byte {
	var filters []any
	for _, not := range f.not {
		filters = append(filters, not)
	}
	for _, and := range f.and {
		filters = append(filters, and)
	}
	filters = append(filters, f.or)

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

func FilterByAttribute(attr string, filters []string) []string {
	fn := func(f string, _ int) (string, bool) {
		pre := attr + ":"
		return strings.TrimPrefix(f, pre), strings.HasPrefix(f, pre)
	}
	return lo.FilterMap(filters, fn)
}
