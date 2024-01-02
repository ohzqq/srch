package srch

import (
	"fmt"

	"github.com/ohzqq/facet"
	"github.com/samber/lo"
)

type Searcher interface {
	Search(...Query) ([]Result, error)
}

type SearchFunc func(...Query) ([]Result, error)

type Query interface {
	fmt.Stringer
}

type Search struct {
	interactive bool
	idx         *facet.Index
	search      Searcher
	results     []Result
	query       string
}

type Opt func(*Search)

func New(s Searcher, opts ...Opt) *Search {
	search := &Search{
		search: s,
	}
	for _, opt := range opts {
		opt(search)
	}
	return search
}

func Interactive(s *Search) {
	s.interactive = true
}

func (s *Search) Get(q ...Query) (*Results, error) {
	if len(q) > 0 {
		s.query = q[0].String()
	}
	var err error
	s.results, err = s.search.Search(q...)
	if err != nil {
		return nil, err
	}

	if s.interactive {
		return s.Choose()
	}

	return s.Results()
}

func (m *Search) Results() (*Results, error) {
	return m.getResults(), nil
}

func (m *Search) getResults(ids ...int) *Results {
	r := &Results{
		Query: m.query,
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

func (s *Search) Choose() (*Results, error) {
	ids, err := Choose(s.results)
	if err != nil {
		return nil, err
	}

	res := s.getResults(ids...)

	return res, nil
}
