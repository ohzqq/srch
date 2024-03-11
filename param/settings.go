package param

import (
	"net/url"
	"slices"
)

type Settings struct {
	params url.Values
}

func NewSettings() *Settings {
	return &Settings{
		params: make(url.Values),
	}
}

func (p Settings) IsFacet(attr string) bool {
	return slices.Contains(p.FacetAttr(), attr)
}

func (p Settings) SrchAttr() []string {
	if p.params.Has(SrchAttr) {
		return p.params[SrchAttr]
	}
	return []string{}
}

func (p Settings) FacetAttr() []string {
	if p.params.Has(FacetAttr) {
		return p.params[FacetAttr]
	}
	return []string{}
}

func (p Settings) SortAttr() []string {
	sort := p.params[SortAttr]
	if sort != nil {
		return sort
	}
	return []string{"title:text"}
}

func (p Settings) HasData() bool {
	return p.params.Has(DataFile) ||
		p.params.Has(DataDir)
}

func (p Settings) GetData() string {
	switch {
	case p.params.Has(DataFile):
		return p.params.Get(DataFile)
	case p.params.Has(DataDir):
		return p.params.Get(DataDir)
	default:
		return ""
	}
}

func (p *Settings) IsFullText() bool {
	return p.params.Has(FullText)
}

func (p *Settings) GetFullText() string {
	return p.params.Get(FullText)
}

func (p Settings) GetUID() string {
	return p.params.Get("uid")
}
