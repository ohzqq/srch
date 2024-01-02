package srch

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
)

type TUI struct {
	*list.Model
}

type item string

func Choose(results []Item) ([]int, error) {
	items := make([]list.Item, len(results))
	for i, r := range results {
		items[i] = item(r.String())
	}

	s := &TUI{}
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

func (m *TUI) Init() tea.Cmd { return nil }

func (m *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
