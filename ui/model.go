package ui

import (
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/srch"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

type TUI struct {
	*list.Model
	*srch.Index
}

type item string

func NewTUI(idx *srch.Index) *TUI {
	return &TUI{
		Index: idx,
	}
}

func Choose(idx *srch.Index) (*srch.Index, error) {
	items := SrcToItems(idx)
	sel, err := NewList(items)
	if err != nil {
		return idx, err
	}

	if len(sel) < 1 {
		return idx, nil
	}

	idx.Data = srch.FilteredItems(idx.Data, lo.ToAnySlice(sel))

	return idx, nil
}

func FilterFacet(facet *srch.Field) string {
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

func SrcToStringSlice(src fuzzy.Source) []string {
	items := make([]string, src.Len())
	for i := 0; i < src.Len(); i++ {
		items[i] = src.String(i)
	}
	return items
}

func StringSliceToItems(src []string) []list.Item {
	items := make([]list.Item, len(src))
	for i, d := range src {
		items[i] = item(d)
	}
	return items
}
