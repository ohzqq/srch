//go:build exclude

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

	idx    *srch.Idx
	params *srch.Params

	data  []map[string]any
	query url.Values

	facetLabels []string
	facets      map[string]*Facet

	*Model
	visible    *srch.Response
	facetMenu  *FacetMenu
	Filters    url.Values
	filters    []any
	facet      string
	Selections *srch.Response
}

func New(idx *srch.Idx) *App {
	tui := newApp()
	tui.idx = idx
	tui.updateVisible(idx.Search(""))
	tui.Model = NewModel(SrcToItems(tui.visible))
	return tui
}

func Browse(q url.Values, data []map[string]any) *App {
	tui := newApp()
	tui.idx, _ = srch.New(q)
	tui.filters = tui.idx.Filters()
	tui.params = srch.ParseParams(q)
	tui.updateVisible(tui.idx.Response())
	tui.Model = NewModel(SrcToItems(tui.visible))
	return tui
}

func newApp() *App {
	return &App{
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
			component := NewFacetMenu(c.visible.FacetLabels())

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

func (c *App) SetFilters(field string, vals []string) {
	if len(vals) == 0 {
		return
	}
	f := srch.NewFilter(field, vals...)
	c.filters = append(c.filters, lo.ToAnySlice(f)...)
	c.idx.SetFilters(c.filters)
	c.updateVisible(c.idx.Filter(""))
}

func (c *App) ClearFilters() {
	c.Filters = make(url.Values)
	c.filters = []any{}
	c.idx.SetFilters(c.filters)
	c.updateVisible(c.idx.Search(""))
}

func (c *App) SetSelections(idx *srch.Response) {
	c.Selections = idx
}

func (c *App) updateVisible(idx *srch.Response) {
	c.visible = idx
	facets := c.visible.Facets()
	if len(facets) > 0 {
		c.facets = make(map[string]*Facet)
		for label, field := range facets {
			c.facets[label] = NewFacet(field)
			//c.setFacet(label)
		}
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
