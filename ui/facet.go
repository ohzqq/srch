package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/srch"
)

type Facet struct {
	*Model
}

type FacetMenu struct {
	*Model
}

type FacetMenuProps struct {
	Items []string
}

type FacetProps struct {
	*srch.Facet
}

func (m *Facet) Init(props FacetProps) tea.Cmd {
	m.UpdateProps(props)
	return nil
}

func (m *Facet) Init(props FacetProps) tea.Cmd {
	m.UpdateProps(props)
	return nil
}
