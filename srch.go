package srch

import (
	"fmt"
	"net/url"

	"github.com/samber/lo"
)

type Searcher interface {
	Search(string) (*Results, error)
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

func (m *Search) Results() (*Results, error) {
	return m.getResults(), nil
}

func (m *Search) getResults(ids ...int) *Results {
	r := &Results{}

	if len(ids) > 0 {
		r.Data = make([]any, len(ids))
		for i, id := range ids {
			r.Data[i] = m.results[id]
		}
		return r
	}
	r.Data = lo.ToAnySlice(m.results)

	return r
}

func (s *Search) Choose() (*Results, error) {
	ids, err := Choose(s.results)
	if err != nil {
		return &Results{}, err
	}

	res := s.getResults(ids...)

	return res, nil
}
