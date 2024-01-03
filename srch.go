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
	search := Search{
		search: s,
	}
	return search
}

func Interactive(s *Index) {
	s.Search.interactive = true
}

func (s *Search) Get(q string) (Search, error) {
	var err error
	s.results, err = s.search.Search(q)
	if err != nil {
		return Search{}, err
	}

	if s.interactive {
		return s.Choose()
	}

	return s.Results()
}

func (m *Search) Results() (Search, error) {
	return m.getResults(), nil
}

func (m *Search) getResults(ids ...int) Search {
	r := Search{
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

func (s *Search) Choose() (Search, error) {
	ids, err := Choose(s.results)
	if err != nil {
		return Search{}, err
	}

	res := s.getResults(ids...)

	return res, nil
}
