package facet

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/viper"
)

const (
	SortByCount = `count`
	SortByAlpha = `alpha`
)

type Facet struct {
	Attribute string         `json:"attribute"`
	Items     []*Item        `json:"items"`
	Count     int            `json:"count"`
	Sep       string         `json:"-"`
	SortBy    string         `json:"-"`
	Order     string         `json:"-"`
	kwIdx     map[string]int `json:"-"`
}

func NewFacet(attr string) *Facet {
	f := &Facet{
		Sep:    "/",
		SortBy: "count",
		Order:  "desc",
	}
	parseAttr(f, attr)
	return f
}

func NewFacets(attrs []string) []*Facet {
	fields := make([]*Facet, len(attrs))
	for i, attr := range attrs {
		fields[i] = NewFacet(attr)
	}
	return fields
}

func (f *Facet) MarshalJSON() ([]byte, error) {
	f.Items = f.Keywords()
	f.Count = f.Len()
	f.Attribute = joinAttr(f)
	//field := make(map[string]any)
	//field["facetValues"] = f.Keywords()
	//if f.Len() < 1 {
	//field["facetValues"] = []any{}
	//}
	//field["attribute"] = joinAttr(f)
	//field["count"] = f.Len()
	return json.Marshal(f)
}

func (f *Facet) Keywords() []*Item {
	return f.SortTokens()
}

func (f *Facet) GetValues() []string {
	vals := make([]string, f.Len())
	for i, token := range f.Items {
		vals[i] = token.Value
	}
	return vals
}

func (f *Facet) FindByLabel(label string) *Item {
	for _, token := range f.Items {
		if token.Label == label {
			return token
		}
	}
	return NewItem(label)
}

func (f *Facet) FindByValue(val string) *Item {
	for _, token := range f.Items {
		if token.Value == val {
			return token
		}
	}
	return NewItem(val)
}

func (f *Facet) FindByIndex(ti ...int) []*Item {
	var tokens []*Item
	for _, tok := range ti {
		if tok < f.Len() {
			tokens = append(tokens, f.Items[tok])
		}
	}
	return tokens
}

func (f *Facet) Add(val any, ids []int) {
	for _, token := range f.Tokenize(val) {
		if f.kwIdx == nil {
			f.kwIdx = make(map[string]int)
		}
		if idx, ok := f.kwIdx[token.Value]; ok {
			f.Items[idx].Add(ids...)
		} else {
			idx = len(f.Items)
			f.kwIdx[token.Value] = idx
			token.Add(ids...)
			f.Items = append(f.Items, token)
		}
	}
}

func (f *Facet) Tokenize(val any) []*Item {
	return KeywordTokenizer(val)
}

func (f *Facet) Search(term string) []*Item {
	matches := fuzzy.FindFrom(term, f)
	tokens := make([]*Item, len(matches))
	for i, match := range matches {
		tokens[i] = f.Items[match.Index]
	}
	return tokens
}

func (f *Facet) Filter(val string) *roaring.Bitmap {
	tokens := f.Find(val)
	bits := make([]*roaring.Bitmap, len(tokens))
	for i, token := range tokens {
		bits[i] = token.Bitmap()
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (f *Facet) Find(val any) []*Item {
	var tokens []*Item
	for _, tok := range f.Tokenize(val) {
		if token, ok := f.kwIdx[tok.Value]; ok {
			tokens = append(tokens, f.Items[token])
		}
	}
	return tokens
}

func (f *Facet) Fuzzy(term string) *roaring.Bitmap {
	matches := fuzzy.FindFrom(term, f)
	bits := make([]*roaring.Bitmap, len(matches))
	for i, match := range matches {
		b := f.Items[match.Index].Bitmap()
		bits[i] = b
	}
	return roaring.ParOr(viper.GetInt("workers"), bits...)
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Facet) Len() int {
	return len(f.Items)
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (f *Facet) String(i int) string {
	return f.Items[i].Label
}

func joinAttr(field *Facet) string {
	attr := field.Attribute
	if field.SortBy != "" {
		attr += ":"
		attr += field.SortBy
	}
	if field.Order != "" {
		attr += ":"
		attr += field.Order
	}
	return attr
}

func parseAttr(field *Facet, attr string) {
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
