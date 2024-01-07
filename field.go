package srch

import (
	"github.com/samber/lo"
)

const (
	Keyword = "keyword"
	Text    = "text"
)

type Field struct {
	Attribute string           `json:"attribute"`
	Items     map[string][]int `json:"items,omitempty"`
	FieldType string
}

func NewField(attr string) *Field {
	return &Field{
		Attribute: attr,
		Items:     make(map[string][]int),
		FieldType: Keyword,
	}
}

func (f *Field) Add(value string, ids ...int) {
	if f.FieldType == Text {
		f.addFullText(value, ids)
		return
	}
	f.addTerm(value, ids)
}

func (f *Field) addFullText(text string, ids []int) {
	for _, token := range Tokenizer(text) {
		f.addTerm(token, ids)
	}
}

func (f *Field) addTerm(term string, ids []int) {
	f.Items[term] = append(f.Items[term], ids...)
}

func (f *Field) Search(text string) []int {
	var r []int
	for _, token := range Tokenizer(text) {
		if ids, ok := f.Items[token]; ok {
			r = append(r, ids...)
		}
	}
	return lo.Uniq(r)
}
