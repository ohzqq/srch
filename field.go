package srch

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	Text       = "text"
	Fuzzy      = "fuzzy"
	OrFacet    = "or"
	AndFacet   = "and"
	NotFacet   = `not`
	FacetField = `facet`
)

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	FieldType string `json:"fieldType"`
	SortBy    string
	Order     string
	items     map[string]*FacetItem `json:"-"`
}

func NewTextField(attr string, params ...*Params) *Field {
	f := NewField(attr, params...)
	f.FieldType = Text
	return f
}

func NewFacet(attr string, params ...*Params) *Field {
	f := NewField(attr, params...)
	f.FieldType = FacetField
	return f
}

func NewField(attr string, params ...*Params) *Field {
	f := &Field{
		Sep:    ".",
		SortBy: "count",
		Order:  "desc",
		items:  make(map[string]*FacetItem),
	}
	parseAttr(f, attr)

	if len(params) > 0 {
		f.SortBy = params[0].SortFacetsBy()
	}

	return f
}

func (f *Field) MarshalJSON() ([]byte, error) {
	field := map[string]any{
		"attribute": f.Attribute,
		"sort_by":   f.SortBy,
		"order":     f.Order,
		"items":     f.Items(),
	}

	d, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (f *Field) Items() []*FacetItem {
	var items []*FacetItem
	for _, k := range f.sortedKeys() {
		items = append(items, f.items[k])
	}
	if f.FieldType == Text {
		return items
	}
	switch f.SortBy {
	case "label":
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
	keys := lo.Keys(f.items)
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
	for _, val := range FacetTokenizer(value) {
		f.addTerm(val, cast.ToIntSlice(ids))
	}
}

func (f *Field) AddFullText(value any, ids any) {
	for _, token := range Tokenizer(cast.ToString(value)) {
		f.addTerm(token, cast.ToIntSlice(ids))
	}
}

func (f *Field) addTerm(item *FacetItem, ids []int) {
	if f.items == nil {
		f.items = make(map[string]*FacetItem)
	}
	if _, ok := f.items[item.Value]; !ok {
		f.items[item.Value] = item
	}
	for _, id := range ids {
		if !f.items[item.Value].bits.ContainsInt(id) {
			f.items[item.Value].bits.AddInt(id)
		}
	}
}

func (f *Field) IsFacet() bool {
	return f.FieldType != Text
}

func (f *Field) ListTokens() []string {
	return lo.Keys(f.items)
}

func (f *Field) Search(text string) *roaring.Bitmap {
	if f.IsFacet() {
		if item, ok := f.items[normalizeText(text)]; ok {
			return item.bits
		}
	}

	var bits []*roaring.Bitmap
	for _, token := range Tokenizer(text) {
		if ids, ok := f.items[token.Value]; ok {
			bits = append(bits, ids.bits)
		}
	}
	return processBitResults(bits, AndFacet)
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
func (f *Field) GetItem(term string) *FacetItem {
	for _, item := range f.Items() {
		if term == item.Label {
			return item
		}
	}
	return &FacetItem{}
}

// ListItems returns a string slice of all item values.
func (f *Field) ListItems() []string {
	var items []string
	for _, item := range f.Items() {
		items = append(items, item.Label)
	}
	return items
}

// FuzzyFindItem fuzzy finds an item's value and returns possible matches.
func (f *Field) FuzzyFindItem(term string) []*FacetItem {
	matches := f.FuzzyMatches(term)
	items := make([]*FacetItem, len(matches))
	for i, match := range matches {
		item := f.Items()[match.Index]
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
	return f.Items()[i].Label
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Field) Len() int {
	return len(f.Items())
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
		f.FieldType == OrFacet ||
		f.FieldType == AndFacet
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
