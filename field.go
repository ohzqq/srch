package srch

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/txt"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
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

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	FieldType string `json:"fieldType"`
	SortBy    string
	Order     string
	tokens    map[string]*txt.Token `json:"-"`
}

func NewField(attr string, params ...*Params) *Field {
	f := &Field{
		Sep:    ".",
		SortBy: "count",
		Order:  "desc",
		tokens: make(map[string]*txt.Token),
	}
	parseAttr(f, attr)

	if len(params) > 0 {
		f.SortBy = params[0].SortFacetsBy()
	}

	return f
}

func NewTextField(attr string, params ...*Params) *Field {
	f := NewField(attr, params...)
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
		"items":     f.Tokens(),
	}

	d, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (f *Field) Tokens() []*txt.Token {
	var items []*txt.Token
	for _, k := range f.sortedKeys() {
		items = append(items, f.tokens[k])
	}
	if f.FieldType == Text {
		return items
	}
	switch f.SortBy {
	case "alpha":
		SortItemsByLabel(items)
	default:
		SortItemsByCount(items)
	}
	if f.Order == "asc" {
		slices.Reverse(items)
	}
	return items
}

func (f *Field) sortedKeys() []string {
	keys := lo.Keys(f.tokens)
	slices.Sort(keys)
	return keys
}

func (f *Field) Add(value any, ids ...any) {
	if f.FieldType == Text {
		f.AddFullText(value, ids)
		return
	}
	f.AddToFacet(value, ids)
}

func (f *Field) AddToFacet(value any, ids any) {
	for _, val := range KeywordAnalyzer(value) {
		f.addTerm(val, cast.ToIntSlice(ids))
	}
}

func (f *Field) AddFullText(value any, ids any) {
	for _, token := range FulltextAnalyzer(cast.ToString(value)) {
		f.addTerm(token, cast.ToIntSlice(ids))
	}
}

func (f *Field) addTerm(item *txt.Token, ids []int) {
	if f.tokens == nil {
		f.tokens = make(map[string]*txt.Token)
	}
	if _, ok := f.tokens[item.Value]; !ok {
		f.tokens[item.Value] = item
	}
	for _, id := range ids {
		if !f.tokens[item.Value].Bitmap().ContainsInt(id) {
			f.tokens[item.Value].Bitmap().AddInt(id)
		}
	}
}

func (f *Field) IsFacet() bool {
	return f.FieldType != Text
}

func (f *Field) ListTokens() []string {
	return lo.Keys(f.tokens)
}

func (f *Field) Search(text string) *roaring.Bitmap {
	if f.IsFacet() {
		if item, ok := f.tokens[normalizeText(text)]; ok {
			return item.Bitmap()
		}
	}

	var bits []*roaring.Bitmap
	for _, token := range FulltextAnalyzer(text) {
		if ids, ok := f.tokens[token.Value]; ok {
			bits = append(bits, ids.Bitmap())
		}
	}
	return processBitResults(bits, And)
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
	for _, item := range f.Tokens() {
		if term == item.Label {
			return item
		}
	}
	return &txt.Token{}
}

// ListItems returns a string slice of all item values.
func (f *Field) ListItems() []string {
	var items []string
	for _, item := range f.Tokens() {
		items = append(items, item.Label)
	}
	return items
}

// FuzzyFindItem fuzzy finds an item's value and returns possible matches.
func (f *Field) FuzzyFindItem(term string) []*txt.Token {
	matches := f.FuzzyMatches(term)
	items := make([]*txt.Token, len(matches))
	for i, match := range matches {
		item := f.Tokens()[match.Index]
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
	return f.Tokens()[i].Label
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Field) Len() int {
	return len(f.Tokens())
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
