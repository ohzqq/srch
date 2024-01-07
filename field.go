package srch

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type FieldType string

const (
	Taxonomy FieldType = "taxonomy"
	Text     FieldType = "text"
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
	return NewField(attr, Text)
}

func NewTaxonomyField(attr string) *Field {
	return NewField(attr, Taxonomy)
}

func (f *Field) Add(value any, ids ...any) {
	if f.FieldType == Text {
		f.addFullText(cast.ToString(value), uint32Slice(ids))
		return
	}
	for _, val := range cast.ToStringSlice(value) {
		f.addTerm(val, uint32Slice(ids))
	}
}

func (f *Field) addFullText(text string, ids []uint32) {
	for _, token := range Tokenizer(text) {
		f.addTerm(token, ids)
	}
}

func (f *Field) addTerm(term string, ids []uint32) {
	if f.Items == nil {
		f.Items = make(map[string]*roaring.Bitmap)
	}
	if _, ok := f.Items[term]; !ok {
		f.Items[term] = roaring.New()
	}
	for _, id := range ids {
		if !f.Items[term].Contains(id) {
			f.Items[term].Add(id)
		}
	}
}

func (f *Field) ListTokens() []string {
	return lo.Keys(f.Items)
}

func (f *Field) Search(text string) *roaring.Bitmap {
	if f.FieldType == Taxonomy {
		if ids, ok := f.Items[text]; ok {
			return ids
		}
	}
	var r []*roaring.Bitmap
	for _, token := range Tokenizer(text) {
		if ids, ok := f.Items[token]; ok {
			r = append(r, ids)
		}
	}
	return roaring.FastOr(r...)
}

func uint32Slice(ids []any) []uint32 {
	bits := make([]uint32, len(ids))
	for i, id := range ids {
		bits[i] = cast.ToUint32(id)
	}
	return bits
}
