package srch

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/txt"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const (
	Text       = "text"
	Fuzzy      = "fuzzy"
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
	tokenz map[string]*txt.Token `json:"-"`
}

func NewField(attr string, params ...*Params) *Field {
	f := &Field{
		Sep:    ".",
		SortBy: "count",
		Order:  "desc",
		Tokens: txt.NewTokens(),
		tokenz: make(map[string]*txt.Token),
	}
	parseAttr(f, attr)

	if len(params) > 0 {
		f.SortBy = params[0].SortFacetsBy()
	}

	return f
}

func NewTextField(attr string, params ...*Params) *Field {
	f := NewField(attr, params...)
	f.Tokens = txt.NewTokens(txt.Fulltext())
	if len(params) > 0 {
		f.FieldType = params[0].GetAnalyzer()
	}
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

func (f *Field) IsFacet() bool {
	return f.FieldType != Text
}

func processBitResults(bits []*roaring.Bitmap, operator string) *roaring.Bitmap {
	switch operator {
	case "and":
		return roaring.ParAnd(viper.GetInt("workers"), bits...)
	default:
		return roaring.ParOr(viper.GetInt("workers"), bits...)
	}
}

// GetItem returns an *FacetItem.
func (f *Field) GetItem(term string) *txt.Token {
	return f.Tokens.GetByLabel(term)
}

// ListItems returns a string slice of all item values.
func (f *Field) ListItems() []string {
	return f.Tokens.GetLabels()
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
func (f *Field) String(i int) string {
	return f.Tokens.Tokens()[i].Label
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Field) Len() int {
	return f.Tokens.Len()
}

func FilterFacets(fields []*Field) []*Field {
	return lo.Filter(fields, filterFacetFields)
}

func FilterTextFields(fields []*Field) []*Field {
	return lo.Filter(fields, filterTextFields)
}

func SearchableFields(fields []*Field) []string {
	f := FilterTextFields(fields)
	return lo.Map(f, mapFieldAttr)
}

func mapFieldAttr(f *Field, _ int) string {
	return f.Attribute
}

func filterTextFields(f *Field, _ int) bool {
	return f.FieldType == Text ||
		f.FieldType == Fuzzy
}

func filterFacetFields(f *Field, _ int) bool {
	return f.FieldType == FacetField ||
		f.FieldType == Or ||
		f.FieldType == And
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
