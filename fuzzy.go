package srch

import (
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

func FuzzyFind(data []map[string]any, q string, fields ...string) *Index {
	idx := New(q)

	if len(data) < 1 {
		return idx
	}

	idx.AddField(GetFieldsFromData(data, fields)...)

	matches := getFuzzyMatches(idx, q)
	fr := getFuzzyResults(data, matches)

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
