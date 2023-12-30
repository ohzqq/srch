package srch

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
)

type Searcher interface {
	Search(...Query) ([]Result, error)
}

type Query interface {
	fmt.Stringer
}

type Result interface {
	fmt.Stringer
}

type Search struct {
	*list.Model
	search  Searcher
	results []Result
}

type result struct {
	value string
}

func New(s Searcher, opts ...Opt) *Search {
	return &Search{
		search: s,
	}
}

func (s *Search) Get(q ...Query) ([]int, error) {
	r, err := s.search(q...)
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, len(r))
	for i, result := range r {
		items[i] = newResult(result)
	}

	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	s.Model = &l

	p := tea.NewProgram(s)
	_, err := p.Run()
	if err != nil {
		return []int{}, err
	}

	return s.ToggledItems(), nil
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
	return m.value
}

func (i *result) Description() string { return "" }
