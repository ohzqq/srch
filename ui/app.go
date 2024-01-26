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

	idx *srch.Index

	data  []map[string]any
	query url.Values

	facetLabels []string
	facets      map[string]*Facet

	*Model
	visible    *srch.Response
	facetMenu  *FacetMenu
	Filters    url.Values
	facet      string
	Selections *srch.Response
}

func New(idx *srch.Index) *App {
	tui := newApp(idx.Params.Values, idx.Data)
	tui.idx = idx
	tui.updateVisible(idx.Search(""))
	tui.Model = NewModel(SrcToItems(tui.visible))
	return tui
}

func Browse(q url.Values, data []map[string]any) *App {
	tui := newApp(q, data)
	idx, _ := srch.New(q)

	tui.updateVisible(idx.Search(""))
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

func (ui *App) Run() (*srch.Response, error) {
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
			component := NewFacetMenu(c.facetLabels)

			return component, component.Init(FacetMenuProps{
				SetFacet: c.SetFacet,
			})
		},
		"facet": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := c.getFacet(c.facet)

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
	c.Filters = srch.ParseQuery(c.Filters, filters)
	c.updateVisible(c.visible.Filter(filters.Encode()))
}

func (c *App) SetSelections(idx *srch.Index) {
	c.Selections = srch.NewResponse(idx)
}

func (c *App) ClearFilters() {
	c.Filters = make(url.Values)
	c.updateVisible(c.idx.Search(""))
}

func (c *App) updateVisible(idx *srch.Response) {
	c.visible = idx
	c.facetLabels = c.visible.FacetLabels()
	if len(c.facetLabels) > 0 {
		//c.facets = make(map[string]*Facet)
		//for _, label := range c.facetLabels {
		//c.setFacet(label)
		//}
	}
}

func (c *App) setFacet(label string) {
	f := c.visible.GetFacet(label)
	c.facets[label] = NewFacet(f)
}

func (c *App) getFacet(label string) *Facet {
	f := c.visible.GetFacet(label)
	return NewFacet(f)
}

func (c *App) Render(w, h int) string {
	return c.mainRouter.Render(w, h)
}

func (ui *App) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return reactea.Destroy
		}
	}
	return ui.mainRouter.Update(msg)
}
