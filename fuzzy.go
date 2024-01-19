package srch

import (
	"net/url"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func FuzzyFind(data []map[string]any, q string, fields ...string) *Index {
	vals := make(url.Values)
	vals.Set("q", q)
	for _, f := range fields {
		vals.Add("field", f)
	}
	idx := NewIndex(vals)

	if len(data) < 1 {
		return idx
	}

	idx.Index(data)

	fr := idx.FuzzyFind(vals.Get("q"))

	return idx.Copy().Index(fr)
}

func getFuzzyMatches(idx *Index, q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, idx)
}

func getFuzzyMatchIndexes(matches fuzzy.Matches) []any {
	return lo.Map(matches, func(m fuzzy.Match, _ int) any {
		return m.Index
	})
}

func getFuzzyResults(data []map[string]any, matches fuzzy.Matches) []map[string]any {
	return FilteredItems(data, getFuzzyMatchIndexes(matches))
}

func getFieldValues(data []map[string]any, fields []string) []string {
	vals := make([]string, len(data))
	for i, item := range data {
		v := make([]string, len(fields))
		for j, f := range fields {
			if str, ok := item[f]; ok {
				v[j] = cast.ToString(str)
			}
		}
		vals[i] = strings.Join(v, " ")
	}
	return vals
}

func GetFieldsFromData(items []map[string]any, names []string) []*Field {
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
			fields = append(fields, NewField(f, Text))
		}
	}
	return fields
}
