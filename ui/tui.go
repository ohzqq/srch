package ui

import (
	"net/url"
	"slices"

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
	facets     map[string]*Model
	facetModel *Model
	query      url.Values
}

func New(idx *srch.Index) *App {
	tui := &App{
		Index:      idx,
		facets:     make(map[string]*Model),
		query:      make(url.Values),
		mainRouter: router.New(),
	}
	//tui.facetModel = tui.FacetMenu()

	tui.Model = newList(idx)

	for _, f := range idx.Facets() {
		tui.facets[f.Attribute] = FacetModel(f)
	}
	return tui
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewModel(SrcToItems(c.Index))

			return component, component.Init(Props{})
		},
		"facetMenu": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewModel(StringSliceToItems(c.listFacets()))

			return component, component.Init(Props{})
		},
	})
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

func (ui *App) Facet(attr string) string {
	var m *Model
	if _, ok := ui.facets[attr]; !ok {
		return ""
	}
	m = ui.facets[attr]
	sel := m.Choose()

	vals := make(url.Values)
	for _, s := range sel {
		vals.Add(attr, m.Items()[s].FilterValue())
	}
	return vals.Encode()
}

func (ui *App) FacetMenu() string {
	facets := ui.listFacets()
	m := ui.facetMenu()
	var sel int
	for _, s := range m.Choose() {
		sel = s
	}
	return facets[sel]
}

func (ui *App) listFacets() []string {
	facets := lo.Keys(ui.facets)
	slices.Sort(facets)
	return facets
}

func (ui *App) facetMenu() *Model {
	facets := ui.listFacets()
	m := NewModel(StringSliceToItems(facets))
	m.SetLimit(1)
	return m
}

func (ui *App) Update(msg tea.Msg) tea.Cmd {
	//var cmds []tea.Cmd
	switch msg := msg.(type) {
	case SelectedFacetMsg:
		return ui.FilterFacetCmd(string(msg))
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

type SelectedFacetMsg string

func (ui *App) SelectFacetCmd() tea.Msg {
	facet := ui.FacetMenu()
	return SelectedFacetMsg(facet)
}

type FilterFacetMsg string

func (ui *App) FilterFacetCmd(attr string) tea.Cmd {
	return func() tea.Msg {
		filter := ui.Facet(attr)
		return FilterFacetMsg(filter)
	}
}
