package txt

import (
	"encoding/json"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/viper"
)

type Tokens struct {
	tokens   map[string]*Token
	analyzer Analyzer
}

type Option func(*Tokens)

func NewTokens(opts ...Option) *Tokens {
	tokens := &Tokens{
		tokens:   make(map[string]*Token),
		analyzer: Simple{},
	}
	for _, opt := range opts {
		opt(tokens)
	}
	return tokens
}

func (t *Tokens) SetAnalyzer(ana Analyzer) *Tokens {
	t.analyzer = ana
	return t
}

func (t *Tokens) Search(val any) *roaring.Bitmap {
	var bits []*roaring.Bitmap

	for _, tok := range t.Tokenize(val) {
		if token, ok := t.tokens[tok.Value]; ok {
			bits = append(bits, token.Bitmap())
		}
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (t *Tokens) Add(val any, ids []int) {
	for _, token := range t.Tokenize(val) {
		t.add(token, ids)
	}
}

func (t *Tokens) Tokenize(val any) []*Token {
	return t.analyzer.Tokenize(val)
}

func (t *Tokens) GetByLabel(label string) *Token {
	for _, token := range t.tokens {
		if token.Label == label {
			return token
		}
	}
	return NewToken(label)
}

func (t *Tokens) add(token *Token, ids []int) {
	if t.tokens == nil {
		t.tokens = make(map[string]*Token)
	}
	if _, ok := t.tokens[token.Value]; !ok {
		t.tokens[token.Value] = token
	}
	t.tokens[token.Value].Add(ids...)
}

func (t *Tokens) Tokens() []*Token {
	tokens := make([]*Token, len(t.tokens))
	for _, t := range t.tokens {
		tokens = append(tokens, t)
	}
	slices.SortStableFunc(tokens, SortByLabelFunc)
	return tokens
}

func (t *Tokens) GetLabels() []string {
	sorted := t.Tokens()
	tokens := make([]string, len(sorted))
	for i, t := range sorted {
		tokens[i] = t.Label
	}
	return tokens
}

func (t *Tokens) GetValues() []string {
	sorted := t.Tokens()
	tokens := make([]string, len(sorted))
	for i, t := range sorted {
		tokens[i] = t.Value
	}
	return tokens
}

func (t *Tokens) Len() int {
	return len(t.tokens)
}

func (t *Tokens) String(i int) string {
	return t.Tokens()[i].Label
}

func (f *Token) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		"value": f.Value,
		"label": f.Label,
		"count": f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func SortItemsByCount(items []*Token) []*Token {
	slices.SortStableFunc(items, SortByCountFunc)
	return items
}

func SortItemsByLabel(items []*Token) []*Token {
	slices.SortStableFunc(items, SortByLabelFunc)
	return items
}

func SortByCountFunc(a *Token, b *Token) int {
	aC := a.Count()
	bC := b.Count()
	switch {
	case aC < bC:
		return 1
	case aC == bC:
		return 0
	default:
		return -1
	}
}

func SortByLabelFunc(a *Token, b *Token) int {
	switch {
	case a.Label < b.Label:
		return 1
	case a.Label == b.Label:
		return 0
	default:
		return -1
	}
}
