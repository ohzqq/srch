package srch

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Results struct {
	Data         []any      `json:"data"`
	SearchFields []string   `json:"search_fields"`
	Facets       []*Facet   `json:"facets"`
	query        Query      `json:"query"`
	Filters      url.Values `json:"filters"`
}

type Item interface {
	fmt.Stringer
}

type DefaultItem struct {
	Value string
}

func (r Results) Search(q ...Queryer) ([]Item, error) {
	var res []Item
	if len(q) > 0 {
		matches := r.FuzzyFind(q[0].String())
		for _, m := range matches {
			res = append(res, &FacetItem{Match: m})
		}
	}
	return res, nil
}

func (r Results) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.Data[i]),
		r.SearchFields,
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}

func (r Results) Len() int {
	return len(r.Data)
}

func (r Results) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, r)
}

func (i *DefaultItem) String() string {
	return i.Value
}
