package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/srch"
)

type Idx struct {
	*Model
}

type IdxProps struct {
	*srch.Index
}

func (m *Idx) Init(props IdxProps) tea.Cmd {
	m.UpdateProps(props)
	return nil
}
