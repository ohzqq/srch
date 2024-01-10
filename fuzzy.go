package srch

import (
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

func FuzzyFind(data []map[string]any, q string, fields ...string) *Index {
	idx := New()

	if len(data) < 1 {
		return idx
	}

	idx.AddField(GetFieldsFromSlice(data, fields)...)

	matches := getFuzzyMatches(idx, q)
	fr := getFuzzyResults(data, matches)

	return New(WithFields(idx.Fields)).Index(fr)
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
	return collectResults(data, getFuzzyMatchIndexes(matches))
}

//func FuzzySearch(data []map[string]any, fields ...string) SearchFunc {
//  return func(q string) []map[string]any {
//    if q == "" {
//      return data
//    }

//    src := GetSearchableFieldValues(data, fields)
//    var res []map[string]any
//    for _, m := range fuzzy.Find(q, src) {
//      res = append(res, data[m.Index])
//    }
//    return res
//  }
//}
