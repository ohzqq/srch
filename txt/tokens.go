package txt

import (
	"github.com/RoaringBitmap/roaring"
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

func (t *Tokens) Filter(val string) *roaring.Bitmap {
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

func (t *Tokens) Add(val any, ids []int) {
	for _, token := range t.Tokenize(val) {
		if t.tokens == nil {
			t.tokens = make(map[string]*Token)
		}
		if _, ok := t.tokens[token.Value]; !ok {
			t.labels = append(t.labels, token.Label)
			t.tokens[token.Value] = token
		}
		t.tokens[token.Value].Add(ids...)
	}
}

func (t *Tokens) Tokenize(val any) []*Token {
	return t.analyzer.Tokenize(val)
}

func (t *Tokens) FindByLabel(label string) *Token {
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
	total := t.Count()
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
		tok := t.FindByLabel(label)
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

func (t *Tokens) Count() int {
	return len(t.tokens)
}
