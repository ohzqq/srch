package srch

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
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

type Result interface {
	fmt.Stringer
}

type Search struct {
	*list.Model
	interactive bool
	idx         *facet.Index
	search      Searcher
	results     []Result
	query       string
}

type Results struct {
	Data   []any          `json:"data"`
	Facets []*facet.Facet `json:"facets"`
	Query  string         `json:"query"`
}

type result struct {
	value string
}

type Opt func(*Search)

func New(s Searcher, opts ...Opt) *Search {
	return &Search{
		search: s,
	}
}

func Interactive(s *Search) {
	s.interactive = true
}

func (s *Search) Get(q ...Query) error {
	if len(q) > 0 {
		s.query = q[0].String()
	}
	var err error
	s.results, err = s.search.Search(q...)
	if err != nil {
		return err
	}
	return nil
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
	items := make([]list.Item, len(s.results))
	for i, result := range s.results {
		items[i] = newResult(result)
	}

	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	s.Model = &l
	s.SetNoLimit()

	p := tea.NewProgram(s)
	_, err := p.Run()
	if err != nil {
		return nil, err
	}

	res := s.getResults(s.ToggledItems()...)

	return res, nil
}

func (m *Search) Init() tea.Cmd { return nil }

func (m *Search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return m, tea.Quit
		}
	}
	l, cmd := m.Model.Update(msg)
	m.Model = &l
	return m, cmd
}

func newResult(r Result) *result {
	return &result{value: r.String()}
}

func (m *result) FilterValue() string {
	return m.value
}

func (i *result) Title() string {
	return i.value
}

func (i *result) Description() string { return "" }
