package ui

import (
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/srch"
	"github.com/samber/lo"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	*srch.Index
	*Model
	og         *srch.Index
	facetModel *Model
	query      url.Values
	facet      string
}

func New(idx *srch.Index) *App {
	tui := &App{
		Index: idx,
		og:    idx.Copy().Index(idx.Data),
		//query:      make(url.Values),
		mainRouter: router.New(),
	}

	tui.Model = newList(idx)

	return tui
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewIdx(c.Index)

			return component, component.Init(IdxProps{
				ClearFilters: c.ClearFilters,
			})
		},
		"filtered": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			idx := c.Index
			if c.query != nil {
				idx = c.Filter(c.query)
			}
			component := NewIdx(idx)

			return component, component.Init(IdxProps{
				ClearFilters: c.ClearFilters,
			})
		},
		"facetMenu": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewFacetMenu(c.FacetLabels())

			return component, component.Init(FacetMenuProps{
				SetFacet: c.SetFacet,
			})
		},
		"facet": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			f, _ := c.GetField(c.facet)
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
	c.query = filters
}

func (c *App) ClearFilters() {
	c.query = nil
	c.Index = c.og
}

func (c *App) Render(w, h int) string {
	return c.mainRouter.Render(w, h)
}

func (ui *App) Choose() (*srch.Index, error) {
	//p := tea.NewProgram(ui)
	p := reactea.NewProgram(ui)
	_, err := p.Run()
	if err != nil {
		return ui.Index, nil
	}

	sel := ui.ToggledItems()
	if len(sel) < 1 {
		return ui.Index, nil
	}

	res := srch.FilteredItems(ui.Index.Data, lo.ToAnySlice(sel))

	return ui.Index.Index(res), nil
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
