package srch

import (
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
