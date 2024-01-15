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
	*srch.Index
	*Model
	facets map[string]*Model
}

type Model struct {
	*list.Model
}

type item string

func NewTUI(idx *srch.Index) *TUI {
	tui := &TUI{
		Index:  idx,
		facets: make(map[string]*Model),
	}

	tui.Model = newList(idx)

	for _, f := range idx.Facets() {
		tui.facets[f.Attribute] = FacetModel(f)
	}
	return tui
}

func (ui *TUI) Choose() (*srch.Index, error) {
	sel := ui.Model.Choose()
	if len(sel) < 1 {
		return ui.Index, nil
	}

	res := srch.FilteredItems(ui.Index.Data, lo.ToAnySlice(sel))

	return ui.Index.Index(res), nil
}

func (ui *TUI) Facet(attr string) url.Values {
	var m *Model
	if _, ok := ui.facets[attr]; !ok {
		return url.Values{}
	}
	m = ui.facets[attr]
	sel := m.Choose()

	vals := make(url.Values)
	for _, s := range sel {
		vals.Add(attr, m.Items()[s].FilterValue())
	}
	return vals
}

func (ui *Model) Choose() []int {
	p := tea.NewProgram(ui)
	_, err := p.Run()
	if err != nil {
		return nil
	}

	return ui.ToggledItems()
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

	res := srch.FilteredItems(idx.Data, lo.ToAnySlice(sel))

	return idx.Index(res), nil
}

func FacetModel(facet *srch.Field) *Model {
	return newList(facet)
}

func newList(src fuzzy.Source) *Model {
	items := SrcToItems(src)
	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	l.SetNoLimit()
	return &Model{
		Model: &l,
	}
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
	s := &Model{}
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

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
