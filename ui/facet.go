package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/srch"
)

type Facet struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FacetProps]
	Model *Model
	*srch.Field
}

type FacetProps struct {
	SetFilters func([]any)
}

func NewFacet(facet *srch.Field) *Facet {
	return &Facet{
		Field: facet,
		Model: NewModel(SrcToItems(facet)),
	}
}

func (m *Facet) Init(props FacetProps) tea.Cmd {
	m.UpdateProps(props)
	return nil
}

func (m *Facet) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.Model.SettingFilter() {
				var filters []any
				for _, s := range m.Model.ToggledItems() {
					f := m.Attribute + ":" + m.Model.Items()[s].FilterValue()
					filters = append(filters, f)
				}
				m.Props().SetFilters(filters)
				reactea.SetCurrentRoute("filtered")
				return nil
			}
		}
	}
	cmds = append(cmds, m.Model.Update(msg))
	return tea.Batch(cmds...)
}

func (m *Facet) Render(w, h int) string {
	return m.Model.Render(w, h)
}
