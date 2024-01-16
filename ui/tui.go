package ui

import (
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/srch"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	*srch.Index
	visible *srch.Index
	*Model
	query      url.Values
	Filters    url.Values
	data       []map[string]any
	facet      string
	Selections *srch.Index
}

func New(idx *srch.Index) *App {
	tui := &App{
		visible:    idx,
		query:      idx.Query,
		data:       idx.Data,
		mainRouter: router.New(),
	}

	tui.Model = newList(idx)

	return tui
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewIdx(c.visible)

			return component, component.Init(IdxProps{
				ClearFilters:  c.ClearFilters,
				SetSelections: c.SetSelections,
			})
		},
		"filtered": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewIdx(c.visible)

			return component, component.Init(IdxProps{
				ClearFilters:  c.ClearFilters,
				SetSelections: c.SetSelections,
			})
		},
		"facetMenu": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewFacetMenu(c.visible.FacetLabels())

			return component, component.Init(FacetMenuProps{
				SetFacet: c.SetFacet,
			})
		},
		"facet": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			f, _ := c.visible.GetField(c.facet)
			component := NewFacet(f)

			return component, component.Init(FacetProps{
				SetFilters: c.SetFilters,
			})
		},
	})
}

func (c *App) SetFacet(label string) {
	c.facet = label
}

func (c *App) SetFilters(filters url.Values) {
	c.Filters = srch.NewQuery(c.Filters, filters)
	c.visible = c.visible.Filter(filters)
	//c.query = srch.NewQuery(c.query, filters)
	//c.visible = srch.New(c.query).Index(data)
}

func (c *App) SetSelections(idx *srch.Index) {
	c.Selections = idx
}

func (c *App) ClearFilters() {
	c.Filters = make(url.Values)
	c.visible = srch.New(c.query).Index(c.data)
}

func (c *App) Render(w, h int) string {
	return c.mainRouter.Render(w, h)
}

func (ui *App) Choose() (*srch.Index, error) {
	//p := tea.NewProgram(ui)
	p := reactea.NewProgram(ui)
	_, err := p.Run()
	if err != nil {
		return ui.visible, err
	}

	return ui.Selections, nil
}

func (ui *App) Update(msg tea.Msg) tea.Cmd {
	//var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		//case "f":
		//reactea.SetCurrentRoute("facetMenu")
		//return ui.NewStatusMessage("facetMenu")
		}
	}
	//cmd := ui.Model.Update(msg)
	//cmds = append(cmds, cmd)
	//ui.Model = ui.(*Model)

	//return tea.Batch(cmds...)
	return ui.mainRouter.Update(msg)
}
