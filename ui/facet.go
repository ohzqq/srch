package ui

import (
	"net/url"

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
	SetFilters func(url.Values)
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
				toggled := m.Model.ToggledItems()
				vals := make([]string, len(toggled))
				for i, s := range toggled {
					vals[i] = m.Model.Items()[s].FilterValue()
				}
				m.Props().SetFilters(m.Attribute, vals)
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
