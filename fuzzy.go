package srch

import (
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

func FuzzyFind(data []map[string]any, q string, fields ...string) *Results {
	idx := New()

	res := NewResults(idx, data)
	if len(data) < 1 {
		return res
	}

	if len(fields) < 1 {
		fields = lo.Keys(data[0])
	}

	for _, t := range fields {
		res.idx.AddField(NewTextField(t))
	}

	matches := getFuzzyMatches(res, q)
	fr := getFuzzyResults(data, matches)

	return NewResults(idx, fr)
}

func getFuzzyMatches(res *Results, q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, res)
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
