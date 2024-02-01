package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/srch"
)

type Idx struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[IdxProps]
	Model *Model
	*srch.Response
	Enter        key.Binding
	FacetMenu    key.Binding
	ClearFilters key.Binding
}

type IdxProps struct {
	ClearFilters  func()
	SetSelections func(*srch.Response)
}

func NewIdx(idx *srch.Response) *Idx {
	m := &Idx{
		Response: idx,
		Model:    NewModel(SrcToItems(idx)),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "return selections"),
		),
		FacetMenu: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "list facets"),
		),
		ClearFilters: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "clear filters"),
		),
	}
	del := list.NewDefaultDelegate()
	keys := []key.Binding{
		m.FacetMenu,
		m.ClearFilters,
	}
	del.ShortHelpFunc = func() []key.Binding {
		return keys
	}
	del.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{keys}
	}
	m.Model.SetDelegate(del)
	return m
}

func (m *Idx) Init(props IdxProps) tea.Cmd {
	m.UpdateProps(props)
	return nil
}

func (m *Idx) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	if reactea.CurrentRoute() == "filtered" {
		cmds = append(cmds, m.Model.NewStatusMessage(m.Get(srch.FacetFilters)))
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.FacetMenu):
			reactea.SetCurrentRoute("facetMenu")
			return nil
		case key.Matches(msg, m.ClearFilters):
			m.Props().ClearFilters()
			reactea.SetCurrentRoute("default")
			return nil
		case key.Matches(msg, m.Enter):
			if !m.Model.SettingFilter() {
				sel := m.Model.ToggledItems()

				if len(sel) < 1 {
					m.Props().SetSelections(m.Index.Response())
					return reactea.Destroy
				}

				res := m.Index.FilterID(sel...)
				//idx, err := srch.New(m.Index.Params.Values())
				//if err != nil {
				//idx = m.Index
				//}
				m.Props().SetSelections(res)
				return reactea.Destroy
			}
		}
	}
	cmds = append(cmds, m.Model.Update(msg))
	return tea.Batch(cmds...)
}

func (m *Idx) Render(w, h int) string {
	return m.Model.Render(w, h)
}
