package srch

import (
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
	"github.com/sahilm/fuzzy"
)

type TUI struct {
	*list.Model
	*Index
}

type item string

func Choose(idx *Index) ([]int, error) {

	items := SrcToItems(idx)

	return NewList(items)
}

func FilterFacet(facet *Facet) string {
	items := SrcToItems(facet)
	sel, err := NewList(items)
	if err != nil {
		return ""
	}
	vals := make(url.Values)
	for _, s := range sel {
		vals.Add(facet.Attribute, items[s].FilterValue())
	}
	return vals.Encode()
}

func NewList(items []list.Item) ([]int, error) {
	s := &TUI{}
	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	s.Model = &l
	s.SetNoLimit()

	p := tea.NewProgram(s)
	_, err := p.Run()
	if err != nil {
		return nil, err
	}

	return s.ToggledItems(), nil
}

func (m *TUI) Init() tea.Cmd { return nil }

func (m *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (i item) FilterValue() string {
	return string(i)
}

func (i item) Title() string {
	return string(i)
}

func (i item) Description() string { return "" }

func SrcToItems(src fuzzy.Source) []list.Item {
	items := make([]list.Item, src.Len())
	for i := 0; i < src.Len(); i++ {
		items[i] = item(src.String(i))
	}
	return items
}
