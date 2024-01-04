package srch

import (
	"fmt"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Item interface {
	fmt.Stringer
}

func NewDefaultItem(val string) *FacetItem {
	return &FacetItem{Value: val}
}

//func (r Results) Search(q string) ([]Item, error) {
//  var res []Item
//  if q == "" {
//    for _, m := range r.Data {
//      item := cast.ToStringMap(m)
//      res = append(res, NewDefaultItem(item["title"].(string)))
//    }
//    return res, nil
//  }
//  matches := r.FuzzyFind(q)
//  for _, m := range matches {
//    res = append(res, &FacetItem{Match: m})
//  }
//  return res, nil
//}

func (r *Index) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.Data[i]),
		r.SearchFields,
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}

func (r *Index) Len() int {
	return len(r.Data)
}

func (r *Index) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, r)
}
