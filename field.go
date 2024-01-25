package srch

import (
	"encoding/json"
	"strings"

	"github.com/ohzqq/srch/txt"
	"github.com/sahilm/fuzzy"
)

const (
	Text       = "text"
	Keyword    = "keyword"
	Or         = "or"
	And        = "and"
	Not        = `not`
	FacetField = `facet`
)

type Analyzer interface {
	Tokenize(string) []*txt.Token
}

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	FieldType string `json:"fieldType"`
	SortBy    string
	Order     string
	*txt.Tokens
}

func NewField(attr string) *Field {
	f := &Field{
		Sep:    ".",
		SortBy: "count",
		Order:  "desc",
		Tokens: txt.NewTokens(txt.Keyword()),
	}
	parseAttr(f, attr)

	return f
}

func (f *Field) MarshalJSON() ([]byte, error) {
	field := map[string]any{
		"attribute": f.Attribute,
		"sort_by":   f.SortBy,
		"order":     f.Order,
		"items":     f.Tokens.Tokens(),
	}

	d, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// FuzzyFindItem fuzzy finds an item's value and returns possible matches.
func (f *Field) FuzzyFindItem(term string) []*txt.Token {
	matches := f.FuzzyMatches(term)
	items := make([]*txt.Token, len(matches))
	for i, match := range matches {
		item := f.Tokens.Tokens()[match.Index]
		item.Match = match
		items[i] = item
	}
	return items
}

// FuzzyMatches returns the fuzzy.Matches of the search.
func (f *Field) FuzzyMatches(term string) fuzzy.Matches {
	return fuzzy.FindFrom(term, f)
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
// Len returns the number of items, to satisfy the fuzzy.Source interface.

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
