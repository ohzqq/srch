package srch

import (
	"github.com/samber/lo"
)

type Searcher interface {
	Search(string) ([]Item, error)
}

type Searchz struct {
	interactive bool
	search      Searcher
	results     []Item
	query       string
}

func NewSearch(s Searcher) Search {
	search := Search{}
	return search
}

func Interactive(s *Index) {
	s.interactive = true
}

func (s *Index) get(q string) (*Index, error) {
	var err error
	s.results, err = s.search.Search(q)
	if err != nil {
		return &Index{}, err
	}

	if s.interactive {
		return s.Choose()
	}

	return s.Results()
}

func (m *Index) Results() (*Index, error) {
	return m.getResults(), nil
}

func (m *Index) getResults(ids ...int) *Index {
	r := &Index{
		//Query: m.query,
	}

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

func (s *Index) Choose() (*Index, error) {
	ids, err := Choose(s.results)
	if err != nil {
		return &Index{}, err
	}

	res := s.getResults(ids...)

	return res, nil
}
