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

	*Model
	visible    *srch.Index
	Filters    url.Values
	data       []map[string]any
	query      url.Values
	facet      string
	Selections *srch.Index
}

func New(idx *srch.Index) *App {
	tui := newApp(idx.Query, idx.Data)
	tui.visible = idx
	tui.Model = NewModel(SrcToItems(tui.visible))
	return tui
}

func Browse(q url.Values, data []map[string]any) *App {
	tui := newApp(q, data)
	tui.visible = srch.New(q).Index(data)
	tui.Model = NewModel(SrcToItems(tui.visible))
	return tui
}

func newApp(q url.Values, data []map[string]any) *App {
	return &App{
		query:      q,
		data:       data,
		mainRouter: router.New(),
	}
}

func (ui *App) Run() (*srch.Index, error) {
	p := reactea.NewProgram(ui)
	_, err := p.Run()
	if err != nil {
		return ui.visible, err
	}

	return ui.Selections, nil
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default":  c.idxComponent,
		"filtered": c.idxComponent,
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

func (c *App) idxComponent(router.Params) (reactea.SomeComponent, tea.Cmd) {
	component := NewIdx(c.visible)

	return component, component.Init(IdxProps{
		ClearFilters:  c.ClearFilters,
		SetSelections: c.SetSelections,
	})
}

func (c *App) SetFacet(label string) {
	c.facet = label
}

func (c *App) SetFilters(filters url.Values) {
	c.Filters = srch.NewQuery(c.Filters, filters)
	c.visible = c.visible.Filter(filters)
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

func (ui *App) Update(msg tea.Msg) tea.Cmd {
	//var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return reactea.Destroy
		}
	}
	return ui.mainRouter.Update(msg)
}
