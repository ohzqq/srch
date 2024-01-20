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
	Not []string
	And []string
	Or  []string
}

type FilterStr string

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

func (f *Filters) Encode() string {
	return f.ToValues().Encode()
}

func ParseFilterJSONString(filters string) (*Filters, error) {

	dec, err := url.QueryUnescape(filters)
	if err != nil {
		return nil, err
	}

	var filter []any
	err = json.Unmarshal([]byte(dec), &filter)
	if err != nil {
		return nil, err
	}

	f := &Filters{}
	for _, v := range filter {
		switch val := v.(type) {
		case string:
			f.And = append(f.And, val)
		case []any:
			f.Or = cast.ToStringSlice(val)
		}
	}
	return f, nil
}

func (f *Filters) String() string {
	var filters []any
	for _, not := range f.Not {
		filters = append(filters, not)
	}
	for _, and := range f.And {
		filters = append(filters, and)
	}
	filters = append(filters, f.Or)

	filter, err := json.Marshal(filters)
	if err != nil {
		filter = []byte{}
	}

	return string(filter)
}

func (f *Filters) ToValues() url.Values {
	vals := make(url.Values)
	vals.Set("facetFilters", f.String())

	return vals
}

func ParseFilters(q any) *Filters {
	filters := &Filters{}

	vals := NewQuery(q)

	if !vals.Has("facetFilters") {
		return filters
	}

	fs := vals.Get("facetFilters")
	for _, filter := range strings.Split(fs, ",") {
		dec, err := url.QueryUnescape(filter)
		if err != nil {
			break
		}
		println(dec)
	}
	return filters
}
