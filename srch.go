package srch

import (
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

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
