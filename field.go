package srch

import (
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
	lex       map[string]*FacetItem `json"-"`
}

func NewField(attr string, ft string) *Field {
	f := &Field{
		Attribute: attr,
		FieldType: ft,
		Sep:       ".",
		SortBy:    "count",
		Order:     "desc",
		lex:       make(map[string]*FacetItem),
	}
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

func (f *Field) FacetItems() []*FacetItem {
	var items []*FacetItem
	for k, bits := range f.lex {
		f.lex[k].Count = len(bits.bits.ToArray())
		items = append(items, f.lex[k])
	}
	return items
}

func (f *Field) addTerm(item *FacetItem, ids []int) {
	if f.lex == nil {
		//f.Items = make(map[string]*roaring.Bitmap)
		f.lex = make(map[string]*FacetItem)
	}
	if _, ok := f.lex[item.Value]; !ok {
		//f.Items[item.Value] = roaring.New()
		f.lex[item.Value] = item
	}
	for _, id := range ids {
		if !f.lex[item.Value].bits.ContainsInt(id) {
			f.lex[item.Value].bits.AddInt(id)
		}
		//if !f.Items[item.Value].ContainsInt(id) {
		//  f.Items[item.Value].AddInt(id)
		//}
	}
}

func (f *Field) ListTokens() []string {
	return lo.Keys(f.lex)
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
		if item, ok := f.lex[normalizeText(text)]; ok {
			return item.bits
		}
	}
	var bits []*roaring.Bitmap
	for _, token := range Tokenizer(text) {
		if ids, ok := f.lex[token.Value]; ok {
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
