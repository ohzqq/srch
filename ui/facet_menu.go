package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type FacetMenu struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FacetMenuProps]
	Model *Model
}

type FacetMenuProps struct {
	SetFacet func(string)
}

func NewFacetMenu(labels []string) *FacetMenu {
	m := NewModel(StringSliceToItems(labels))
	m.SetLimit(1)
	return &FacetMenu{
		Model: m,
	}
}

func (m *FacetMenu) Init(props FacetMenuProps) tea.Cmd {
	m.UpdateProps(props)
	return nil
}

func (m *FacetMenu) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.Model.SettingFilter() {
				item := m.Model.SelectedItem()
				m.Props().SetFacet(item.FilterValue())
				reactea.SetCurrentRoute("facet")
				return nil
			}
		}
	}
	cmds = append(cmds, m.Model.Update(msg))
	return tea.Batch(cmds...)
}

func (m *FacetMenu) Render(w, h int) string {
	return m.Model.Render(w, h)
}
