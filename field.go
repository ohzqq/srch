package srch

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type FieldType string

const (
	FacetField FieldType = "facet"
	Text       FieldType = "text"
)

type Field struct {
	Attribute string    `json:"attribute"`
	Operator  string    `json:"operator,omitempty"`
	Sep       string    `json:"-"`
	FieldType FieldType `json:"fieldType"`
	Items     map[string]*roaring.Bitmap
}

func NewField(attr string, ft FieldType) *Field {
	return &Field{
		Attribute: attr,
		FieldType: ft,
		Items:     make(map[string]*roaring.Bitmap),
	}
}

func NewTextField(attr string) *Field {
	f := NewField(attr, Text)
	f.Operator = "and"
	return f
}

func NewTaxonomyField(attr string) *Field {
	f := NewField(attr, FacetField)
	f.Operator = "or"
	return f
}

func (f *Field) Add(value any, ids ...any) {
	if f.FieldType == Text {
		f.addFullText(cast.ToString(value), cast.ToIntSlice(ids))
		return
	}
	for _, val := range cast.ToStringSlice(value) {
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

// GetConfig returns a map of a Facet's config.
func (f *Field) GetConfig() map[string]any {
	return map[string]any{
		"attribute": f.Attribute,
		"operator":  f.Operator,
		"fieldType": f.FieldType,
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
	return f.processResults(bits)
}

func (f *Field) Search(text string) *roaring.Bitmap {
	if f.FieldType == FacetField {
		if ids, ok := f.Items[text]; ok {
			return ids
		}
	}
	var bits []*roaring.Bitmap
	for _, token := range Tokenizer(text) {
		if ids, ok := f.Items[token]; ok {
			bits = append(bits, ids)
		}
	}
	return f.processResults(bits)
}

func (f *Field) processResults(bits []*roaring.Bitmap) *roaring.Bitmap {
	switch f.Operator {
	case "and":
		return roaring.ParAnd(viper.GetInt("workers"), bits...)
	default:
		return roaring.ParOr(viper.GetInt("workers"), bits...)
	}
}
