package ui

import (
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/srch"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

type Model struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	*list.Model
}

type item string

type Props struct {
}

func NewModel(items []list.Item) *Model {
	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	l.SetNoLimit()
	return &Model{
		Model: &l,
	}
}

func (ui *Model) Choose() []int {
	//p := tea.NewProgram(ui)
	//_, err := p.Run()
	//if err != nil {
	//  return nil
	//}

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
	return NewModel(SrcToItems(src))
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

	//p := tea.NewProgram(s)
	//_, err := p.Run()
	//if err != nil {
	//return nil, err
	//}

	return s.ToggledItems(), nil
}

func (m *Model) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	return nil
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "f":
			reactea.SetCurrentRoute("facetMenu")
			//println("facetMenu")
			return m.NewStatusMessage("facetMenu")
		case "enter":
			if !m.Model.SettingFilter() {
				if !m.MultiSelectable() {
					m.ToggleItem()
				}
				return tea.Quit
			}
		}
	}
	l, cmd := m.Model.Update(msg)
	cmds = append(cmds, cmd)
	m.Model = &l
	return tea.Batch(cmds...)
}

func (m *Model) Render(w, h int) string {
	return m.Model.View()
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
