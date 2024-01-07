package srch

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type FieldType string

const (
	Taxonomy FieldType = "keyword"
	Text     FieldType = "text"
)

type Field struct {
	Attribute string              `json:"attribute"`
	Items     map[string][]uint32 `json:"items,omitempty"`
	Operator  string              `json:"operator,omitempty"`
	Sep       string              `json:"-"`
	FieldType FieldType           `json:"fieldType"`
}

func NewField(attr string, ft FieldType) *Field {
	return &Field{
		Attribute: attr,
		Items:     make(map[string][]uint32),
		FieldType: ft,
	}
}

func NewTextField(attr string) *Field {
	return NewField(attr, Text)
}

func NewTaxonomyField(attr string) *Field {
	return NewField(attr, Taxonomy)
}

func (f *Field) Add(value string, ids ...any) {
	if f.FieldType == Text {
		f.addFullText(value, uint32Slice(ids))
		return
	}
	f.addTerm(value, uint32Slice(ids))
}

func (f *Field) addFullText(text string, ids []uint32) {
	for _, token := range Tokenizer(text) {
		f.addTerm(token, ids)
	}
}

func (f *Field) addTerm(term string, ids []uint32) {
	if f.Items == nil {
		f.Items = make(map[string][]uint32)
	}
	f.Items[term] = append(f.Items[term], ids...)
}

func (f *Field) Search(text string) []uint32 {
	var r []uint32
	for _, token := range Tokenizer(text) {
		if ids, ok := f.Items[token]; ok {
			r = append(r, ids...)
		}
	}
	return lo.Uniq(r)
}

// Bitmap returns a *roaring.Bitmap of slice indices for a FacetItem.
func (f *Field) BitmapOf(term string) *roaring.Bitmap {
	if items, ok := f.Items[term]; ok {
		return roaring.BitmapOf(items...)
	}
	return nil
}

func uint32Slice(ids []any) []uint32 {
	bits := make([]uint32, len(ids))
	for i, id := range ids {
		bits[i] = cast.ToUint32(id)
	}
	return bits
}
