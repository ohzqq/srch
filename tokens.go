package srch

import (
	"github.com/ohzqq/srch/txt"
	"github.com/spf13/cast"
)

type Tokens struct {
	tokens   map[string]*Token
	Tokens   []string
	analyzer Analyzer
	ana      *txt.Analyzer
}

func NewTokens() *Tokens {
	tokens := &Tokens{
		tokens:   make(map[string]*Token),
		analyzer: Simple{},
		ana:      txt.Keywords(),
	}
	return tokens
}

func (t *Tokens) SetAnalyzer(ana Analyzer) *Tokens {
	t.analyzer = ana
	return t
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
			t.Tokens = append(t.Tokens, token.Label)
			t.tokens[token.Value] = token
		}
		t.tokens[token.Value].Add(ids...)
	}
}

func (t *Tokens) Tokenize(val any) []*Token {
	tokens, _ := t.ana.Tokenize(cast.ToString(val))
	toks := make([]*Token, len(tokens))
	for i, t := range tokens {
		toks[i] = newTok(t)
	}
	return toks
}

func (t *Tokens) FindByLabel(label string) *Token {
	for _, token := range t.tokens {
		if token.Label == label {
			return token
		}
	}
	return NewToken(label, label)
}

func (t *Tokens) Count() int {
	return len(t.tokens)
}
