package srch

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/ohzqq/srch/txt"
)

const (
	Or          = "or"
	And         = "and"
	Not         = `not`
	FacetField  = `facet`
	SortByCount = `count`
	SortByAlpha = `alpha`
)

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	SortBy    string
	Order     string
	*txt.Tokens
}

func NewField(attr string) *Field {
	f := &Field{
		Sep:    ".",
		SortBy: "count",
		Order:  "desc",
		Tokens: txt.NewTokens(),
	}
	parseAttr(f, attr)
	return f
}

func (f *Field) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(f.Tokens)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (f *Field) GetTokens() []*txt.Token {
	return f.Tokens.Tokens()
}

func (f *Field) Find(kw string) []*txt.Token {
	return f.Tokens.Find(kw)
}

func parseAttr(field *Field, attr string) {
	i := 0
	for attr != "" {
		var a string
		a, attr, _ = strings.Cut(attr, ":")
		if a == "" {
			continue
		}
		switch i {
		case 0:
			field.Attribute = a
		case 1:
			field.SortBy = a
		case 2:
			field.Order = a
		}
		i++
	}
}

func SortItemsByCount(items []*txt.Token) []*txt.Token {
	slices.SortStableFunc(items, SortByCountFunc)
	return items
}

func SortItemsByLabel(items []*txt.Token) []*txt.Token {
	slices.SortStableFunc(items, SortByLabelFunc)
	return items
}

func SortByCountFunc(a *txt.Token, b *txt.Token) int {
	aC := a.Count()
	bC := b.Count()
	switch {
	case aC < bC:
		return 1
	case aC == bC:
		return 0
	default:
		return -1
	}
}

func SortByLabelFunc(a *txt.Token, b *txt.Token) int {
	switch {
	case a.Label > b.Label:
		return 1
	case a.Label == b.Label:
		return 0
	default:
		return -1
	}
}
