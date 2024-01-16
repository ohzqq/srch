package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/srch"
)

type Idx struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[IdxProps]
	Model *Model
	*srch.Index
}

type IdxProps struct {
	ClearFilters func()
}

func NewIdx(idx *srch.Index) *Idx {
	return &Idx{
		Index: idx,
		Model: NewModel(SrcToItems(idx)),
	}
}

func (m *Idx) Init(props IdxProps) tea.Cmd {
	m.UpdateProps(props)
	return nil
}

func (m *Idx) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.Model.SettingFilter() {
			}
		case "f":
			reactea.SetCurrentRoute("facetMenu")
			//println("facetMenu")
			return m.Model.NewStatusMessage("facetMenu")
		case "c":
			m.Props().ClearFilters()
			reactea.SetCurrentRoute("default")
			return nil
		}
	}
	cmds = append(cmds, m.Model.Update(msg))
	return tea.Batch(cmds...)
}

func (m *Idx) Render(w, h int) string {
	return m.Model.Render(w, h)
}
