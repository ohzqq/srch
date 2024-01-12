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
	Attribute string                     `json:"attribute"`
	Operator  string                     `json:"operator,omitempty"`
	Sep       string                     `json:"-"`
	FieldType string                     `json:"fieldType"`
	Items     map[string]*roaring.Bitmap `json:"-"`
}

func NewField(attr string, ft string) *Field {
	f := &Field{
		Attribute: attr,
		FieldType: ft,
		Sep:       ".",
		Items:     make(map[string]*roaring.Bitmap),
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

func (f *Field) addTerm(term string, ids []int) {
	if f.Items == nil {
		f.Items = make(map[string]*roaring.Bitmap)
	}
	if _, ok := f.Items[term]; !ok {
		f.Items[term] = roaring.New()
	}
	for _, id := range ids {
		if !f.Items[term].ContainsInt(id) {
			f.Items[term].AddInt(id)
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
		if ids, ok := f.Items[normalizeText(text)]; ok {
			return ids
		}
	}
	var bits []*roaring.Bitmap
	for _, token := range Tokenizer(text) {
		if ids, ok := f.Items[token]; ok {
			bits = append(bits, ids)
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
	return f.FieldType == FacetField
}
