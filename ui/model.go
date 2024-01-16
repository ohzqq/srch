package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/bubbles/list"
	"github.com/sahilm/fuzzy"
)

type Model struct {
	reactea.BasicComponent

	*list.Model
}

type item string

func NewModel(items []list.Item) *Model {
	return &Model{
		Model: newListModel(items),
	}
}

func newListModel(items []list.Item) *list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	l.SetNoLimit()
	return &l
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	l, cmd := m.Model.Update(msg)
	m.Model = &l
	return cmd
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
