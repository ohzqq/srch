package srch

import (
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	FacetField = "facet"
	Text       = "text"
	OrFacet    = "or"
	AndFacet   = "and"
)

type Field struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator,omitempty"`
	Sep       string `json:"-"`
	FieldType string `json:"fieldType"`
	SortBy    string
	Order     string
	Items     map[string]*FacetItem `json:"-"`
}

func NewField(attr string, ft string) *Field {
	f := &Field{
		FieldType: ft,
		Sep:       ".",
		SortBy:    "count",
		Order:     "desc",
		Items:     make(map[string]*FacetItem),
	}
	parseAttr(f, attr)
	switch ft {
	case OrFacet:
		f.Operator = "or"
	case AndFacet, Text, FacetField:
		f.Operator = "and"
	}
	return f
}

func CopyField(field *Field) *Field {
	f := NewField(field.Attribute, field.FieldType)
	f.Sep = field.Sep
	f.Operator = field.Operator
	return f
}

func NewTextField(attr string) *Field {
	f := NewField(attr, Text)
	f.Operator = "and"
	return f
}

func NewTextFields(names []string) []*Field {
	fields := make([]*Field, len(names))
	for i, f := range names {
		fields[i] = NewTextField(f)
	}
	return fields
}

func NewFacets(names []string) []*Field {
	fields := make([]*Field, len(names))
	for i, f := range names {
		fields[i] = NewFacetField(f)
	}
	return fields
}

func NewFacetField(attr string) *Field {
	f := NewField(attr, FacetField)
	f.Operator = "or"
	return f
}

func (f *Field) Add(value any, ids ...any) {
	if f.FieldType == Text {
		f.addFullText(cast.ToString(value), cast.ToIntSlice(ids))
		return
	}
	for _, val := range FacetTokenizer(value) {
		f.addTerm(val, cast.ToIntSlice(ids))
	}
}

func (f *Field) addFullText(text string, ids []int) {
	for _, token := range Tokenizer(text) {
		f.addTerm(token, ids)
	}
}

func (f *Field) ItemsWithCount() []*FacetItem {
	var items []*FacetItem
	for k, item := range f.Items {
		f.Items[k].Count = len(item.bits.ToArray())
		items = append(items, f.Items[k])
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

func (f *Field) addTerm(item *FacetItem, ids []int) {
	if f.Items == nil {
		f.Items = make(map[string]*FacetItem)
	}
	if _, ok := f.Items[item.Value]; !ok {
		f.Items[item.Value] = item
	}
	for _, id := range ids {
		if !f.Items[item.Value].bits.ContainsInt(id) {
			f.Items[item.Value].bits.AddInt(id)
		}
	}
}

func (f *Field) ListTokens() []string {
	return lo.Keys(f.Items)
}

// Filter applies the listed filters to the facet.
func (f *Field) Filter(filters ...string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, filter := range filters {
		bits = append(bits, f.Search(filter))
	}
	return processBitResults(bits, f.Operator)
}

func (f *Field) Search(text string) *roaring.Bitmap {
	if f.FieldType == FacetField {
		if item, ok := f.Items[normalizeText(text)]; ok {
			return item.bits
		}
	}
	var bits []*roaring.Bitmap
	for _, token := range Tokenizer(text) {
		if ids, ok := f.Items[token.Value]; ok {
			bits = append(bits, ids.bits)
		}
	}
	return processBitResults(bits, f.Operator)
}

func processBitResults(bits []*roaring.Bitmap, operator string) *roaring.Bitmap {
	switch operator {
	case "and":
		return roaring.ParAnd(viper.GetInt("workers"), bits...)
	default:
		return roaring.ParOr(viper.GetInt("workers"), bits...)
	}
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
	return f.FieldType == Text
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
