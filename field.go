package srch

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/ohzqq/srch/txt"
)

const (
	Text       = "text"
	Keyword    = "keyword"
	Or         = "or"
	And        = "and"
	Not        = `not`
	FacetField = `facet`
)

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	SortBy    string
	Order     string
	*txt.Tokens
}

func NewField(attr string, opts ...txt.Option) *Field {
	f := &Field{
		Sep:    ".",
		SortBy: "count",
		Order:  "desc",
		Tokens: txt.NewTokens(opts...),
	}
	parseAttr(f, attr)

	return f
}

func (f *Field) MarshalJSON() ([]byte, error) {
	field := map[string]any{
		"attribute": f.Attribute,
		"sort_by":   f.SortBy,
		"order":     f.Order,
		"items":     f.GetTokens(),
	}

	d, err := json.Marshal(field)
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
