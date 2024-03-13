package srch

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/txt"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	Or               = "OR"
	And              = "AND"
	Not              = `NOT`
	AndNot           = `AND NOT`
	OrNot            = `OR NOT`
	FacetField       = `facet`
	SortByCount      = `count`
	SortByAlpha      = `alpha`
	StandardAnalyzer = `standard`
)

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	SortBy    string
	Order     string
	tokens    map[string]*Token
	Tokens    []string
	ana       *txt.Analyzer
}

func NewField(attr string) *Field {
	f := &Field{
		tokens: make(map[string]*Token),
		Sep:    ".",
		ana:    txt.Keywords(),
	}
	parseAttr(f, attr)
	return f
}

func (t *Field) Add(val any, ids []int) {
	for _, token := range t.Tokenize(val) {
		if t.tokens == nil {
			t.tokens = make(map[string]*Token)
		}
		if _, ok := t.tokens[token.Value]; !ok {
			t.Tokens = append(t.Tokens, token.Label)
			t.tokens[token.Value] = token
		}
		t.tokens[token.Value].Add(ids...)
	}
}

func (t *Field) Find(val any) []*Token {
	var tokens []*Token
	for _, tok := range t.Tokenize(val) {
		if token, ok := t.tokens[tok.Value]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func (t *Field) Search(term string) []*Token {
	matches := fuzzy.FindFrom(term, t)
	tokens := make([]*Token, len(matches))
	all := t.GetTokens()
	for i, match := range matches {
		tokens[i] = all[match.Index]
	}
	return tokens
}

func (t *Field) Filter(val string) *roaring.Bitmap {
	tokens := t.Find(val)
	bits := make([]*roaring.Bitmap, len(tokens))
	for i, token := range tokens {
		bits[i] = token.Bitmap()
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (t *Field) Fuzzy(term string) *roaring.Bitmap {
	matches := t.Search(term)
	bits := make([]*roaring.Bitmap, len(matches))
	for i, match := range matches {
		bits[i] = match.Bitmap()
	}
	return roaring.ParOr(viper.GetInt("workers"), bits...)
}

func (t *Field) FindByLabel(label string) *Token {
	for _, token := range t.tokens {
		if token.Label == label {
			return token
		}
	}
	return NewToken(label, label)
}

func (t *Field) FindByIndex(ti ...int) []*Token {
	var tokens []*Token
	toks := t.GetTokens()
	total := t.Count()
	for _, tok := range ti {
		if tok < total {
			tokens = append(tokens, toks[tok])
		}
	}
	return tokens
}

func (t *Field) SortTokens() []*Token {
	tokens := t.GetTokens()

	switch t.SortBy {
	case SortByAlpha:
		if t.Order == "" {
			t.Order = "asc"
		}
		SortTokensByAlpha(tokens)
	default:
		if t.Order == "" {
			t.Order = "desc"
		}
		SortTokensByCount(tokens)
	}

	if t.Order == "desc" {
		slices.Reverse(tokens)
	}

	return tokens
}

func (f *Field) SetAnalyzer(ana *txt.Analyzer) *Field {
	f.ana = ana
	return f
}

func (f *Field) Keywords() *Field {
	f.ana = txt.Keywords()
	return f
}

func (f *Field) Normalize() *Field {
	f.ana = txt.NewNormalizer()
	return f
}

func (t *Field) Tokenize(val any) []*Token {
	tokens, _ := t.ana.Tokenize(cast.ToString(val))
	toks := make([]*Token, len(tokens))
	for i, t := range tokens {
		toks[i] = newTok(t)
	}
	return toks
}

func (t *Field) Count() int {
	return len(t.tokens)
}

func (f *Field) MarshalJSON() ([]byte, error) {
	tokens := make(map[string]int)
	for _, label := range f.Tokens {
		token := f.FindByLabel(label)
		tokens[label] = token.Count()
	}
	d, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (t *Field) GetTokens() []*Token {
	var tokens []*Token
	for _, label := range t.Tokens {
		tok := t.FindByLabel(label)
		tokens = append(tokens, tok)
	}
	return tokens
}

func GetFieldItems(data []map[string]any, field *Field) []map[string]any {
	field.SortBy = SortByAlpha
	tokens := field.SortTokens()

	items := make([]map[string]any, len(tokens))
	for i, token := range tokens {
		items[i] = map[string]any{
			"attribute": field.Attribute,
			"value":     token.Value,
			"label":     token.Label,
			"count":     token.Count(),
			"hits":      ItemsByBitmap(data, token.Bitmap()),
		}
	}
	return items
}

func (t *Field) GetValues() []string {
	sorted := t.GetTokens()
	tokens := make([]string, len(sorted))
	for i, t := range sorted {
		tokens[i] = t.Value
	}
	return tokens
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (t *Field) Len() int {
	return t.Count()
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (t *Field) String(i int) string {
	return t.Tokens[i]
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

func SortTokensByCount(items []*Token) []*Token {
	slices.SortStableFunc(items, SortByCountFunc)
	return items
}

func SortTokensByAlpha(items []*Token) []*Token {
	slices.SortStableFunc(items, SortByAlphaFunc)
	return items
}

func SortByCountFunc(a *Token, b *Token) int {
	aC := a.Count()
	bC := b.Count()
	switch {
	case aC > bC:
		return 1
	case aC == bC:
		return 0
	default:
		return -1
	}
}

func SortByAlphaFunc(a *Token, b *Token) int {
	aL := strings.ToLower(a.Label)
	bL := strings.ToLower(b.Label)
	switch {
	case aL > bL:
		return 1
	case aL == bL:
		return 0
	default:
		return -1
	}
}
