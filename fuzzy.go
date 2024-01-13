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

	idx.AddField(GetFieldsFromSlice(data, fields)...)

	matches := getFuzzyMatches(idx, q)
	fr := getFuzzyResults(data, matches)

	return idx.Copy().Index(fr)
}

func getFuzzyMatches(idx *Index, q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, idx)
}

func getFuzzyMatchIndexes(matches fuzzy.Matches) []int {
	return lo.Map(matches, func(m fuzzy.Match, _ int) int {
		return m.Index
	})
}

func getFuzzyResults(data []map[string]any, matches fuzzy.Matches) []map[string]any {
	return FilterDataByID(data, getFuzzyMatchIndexes(matches))
}
