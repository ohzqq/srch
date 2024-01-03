package srch

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Searcher interface {
	Search([]any, string) ([]Item, error)
}

type Search struct {
	SearchFields []string   `json:"search_fields"`
	query        Query      `json:"query"`
	Facets       []*Facet   `json:"facets"`
	Filters      url.Values `json:"filters"`
	interactive  bool
	results      []Item
	data         []any
}

type Item interface {
	fmt.Stringer
}

func NewSearch() *Search {
	search := Search{
		SearchFields: []string{"title"},
	}
	return &search
}

func Interactive(s *Index) {
	s.Search.interactive = true
}

func NewDefaultItem(val string) *FacetItem {
	return &FacetItem{Value: val}
}

func (r Search) Search(data []any, q string) ([]Item, error) {
	r.data = data
	var res []Item
	if q == "" {
		for _, m := range data {
			item := cast.ToStringMap(m)
			res = append(res, NewDefaultItem(item["title"].(string)))
		}
		return res, nil
	}
	matches := r.FuzzyFind(q)
	for _, m := range matches {
		res = append(res, &FacetItem{Match: m})
	}
	return res, nil
}

func (r Search) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.data[i]),
		r.SearchFields,
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}

func (r Search) Len() int {
	return len(r.data)
}

func (r Search) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, r)
}

func (m *Search) Results() (Search, error) {
	return m.getResults(), nil
}

func (m *Search) getResults(ids ...int) Search {
	r := Search{
		//Query: m.query,
	}

	if len(ids) > 0 {
		r.data = make([]any, len(ids))
		for i, id := range ids {
			r.data[i] = m.results[id]
		}
		return r
	}
	r.data = lo.ToAnySlice(m.results)

	return r
}

func (s *Search) Choose() (Search, error) {
	ids, err := Choose(s.results)
	if err != nil {
		return Search{}, err
	}

	res := s.getResults(ids...)

	return res, nil
}
