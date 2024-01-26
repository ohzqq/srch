package txt

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/viper"
)

type Tokens struct {
	tokens   map[string]*Token
	labels   []string
	analyzer Analyzer
}

func NewTokens() *Tokens {
	tokens := &Tokens{
		tokens:   make(map[string]*Token),
		analyzer: Simple{},
	}
	return tokens
}

func (t *Tokens) SetAnalyzer(ana Analyzer) *Tokens {
	t.analyzer = ana
	return t
}

func (t *Tokens) Filter(val any) *roaring.Bitmap {
	tokens := t.Find(val)
	bits := make([]*roaring.Bitmap, len(tokens))
	for i, token := range tokens {
		bits[i] = token.Bitmap()
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (t *Tokens) Find(val any) []*Token {
	var tokens []*Token
	for _, tok := range t.Tokenize(val) {
		if token, ok := t.tokens[tok.Value]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func (t *Tokens) Search(term string) []*Token {
	//tokens := t.Fuzzy(term)
	matches := fuzzy.FindFrom(term, t)
	tokens := make([]*Token, len(matches))
	all := t.Tokens()
	for i, match := range matches {
		tokens[i] = all[match.Index]
	}
	return tokens
}

func (t *Tokens) Fuzzy(term string) *roaring.Bitmap {
	matches := fuzzy.FindFrom(term, t)
	all := t.Tokens()
	bits := make([]*roaring.Bitmap, len(matches))
	for i, match := range matches {
		bits[i] = all[match.Index].Bitmap()
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (t *Tokens) Add(val any, ids []int) {
	for _, token := range t.Tokenize(val) {
		t.add(token, ids)
	}
}

func (t *Tokens) add(token *Token, ids []int) {
	if t.tokens == nil {
		t.tokens = make(map[string]*Token)
	}
	if _, ok := t.tokens[token.Value]; !ok {
		t.labels = append(t.labels, token.Label)
		t.tokens[token.Value] = token
	}
	t.tokens[token.Value].Add(ids...)
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

func (t *Tokens) FindByIndex(ti ...int) []*Token {
	var tokens []*Token
	toks := t.Tokens()
	total := t.Len()
	for _, tok := range ti {
		if tok < total {
			tokens = append(tokens, toks[tok])
		}
	}
	return tokens
}

func (t *Tokens) Tokens() []*Token {
	var tokens []*Token
	for _, label := range t.labels {
		tok := t.GetByLabel(label)
		tokens = append(tokens, tok)
	}
	return tokens
}

func (t *Tokens) GetLabels() []string {
	return t.labels
}

func (t *Tokens) GetValues() []string {
	sorted := t.Tokens()
	tokens := make([]string, len(sorted))
	for i, t := range sorted {
		tokens[i] = t.Value
	}
	return tokens
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (t *Tokens) Len() int {
	return len(t.tokens)
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (t *Tokens) String(i int) string {
	return t.labels[i]
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
