package srch

import (
	"net/url"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

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

func collectResults(d []map[string]any, ids []int) []map[string]any {
	if len(ids) > 0 {
		data := make([]map[string]any, len(ids))
		for i, id := range ids {
			data[i] = d[id]
		}
		return data
	}
	return d
}

// ParseValues takes an interface{} and returns a url.Values.
func ParseValues(f any) (map[string][]string, error) {
	filters := make(map[string][]string)
	var err error
	switch val := f.(type) {
	case url.Values:
		return val, nil
	case []byte:
		return ParseQueryBytes(val)
	case string:
		return ParseQueryString(val)
	default:
		filters, err = cast.ToStringMapStringSliceE(val)
		if err != nil {
			return nil, err
		}
	}
	return filters, nil
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (map[string][]string, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, err
}

// ParseQueryBytes parses a byte slice to url.Values.
func ParseQueryBytes(val []byte) (map[string][]string, error) {
	filters, err := cast.ToStringMapStringSliceE(string(val))
	if err != nil {
		return nil, err
	}
	return filters, err
}
