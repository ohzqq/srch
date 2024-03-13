package txt

import (
	"errors"
	"fmt"
	"slices"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

type Token struct {
	Value string `json:"value"`
	Label string `json:"label"`
	fuzzy.Match
}

type Tokens []*Token

var (
	FieldsFuncErr = errors.New("strings.FieldsFunc returned an empty slice or the string was empty")
	EmptyStrErr   = errors.New("empty string")
	NoMatchErr    = errors.New(`no matches found`)
)

func NewToken(label, val string) *Token {
	return &Token{
		Value: val,
		Label: label,
	}
}

func (toks Tokens) Find(q string) (Tokens, error) {
	var tokens Tokens
	for i, tok := range toks {
		if tok.Value == q {
			tok.Match = newMatch(tok.Value, i)
			tokens = append(tokens, tok)
		} else if tok.Label == q {
			tok.Match = newMatch(tok.Label, i)
			tokens = append(tokens, tok)
		}
	}
	if tokens.Len() > 0 {
		return tokens, nil
	}

	return nil, fmt.Errorf("%w for query '%s'\n", NoMatchErr, q)
}

func (toks Tokens) Without(sw Tokens) Tokens {
	return lo.Without(toks, sw...)
}

func (toks Tokens) Search(q string) (Tokens, error) {
	var tokens Tokens
	for _, m := range fuzzy.FindFrom(q, toks) {
		tok := toks[m.Index]
		tok.Match = m
		tok.Match.Str = tok.Label
		tokens = append(tokens, tok)
	}

	if tokens.Len() > 0 {
		return tokens, nil
	}
	return nil, fmt.Errorf("%w for query '%s'\n", NoMatchErr, q)
}

func (toks Tokens) FindByLabel(label string) (*Token, error) {
	for i, token := range toks {
		if token.Label == label {
			token.Match = newMatch(token.Label, i)
			return token, nil
		}
	}
	return nil, fmt.Errorf("%w for label '%s'\n", NoMatchErr, label)
}

func (toks Tokens) FindByValue(val string) (*Token, error) {
	for i, token := range toks {
		if token.Value == val {
			token.Match = newMatch(token.Value, i)
			return token, nil
		}
	}
	return nil, fmt.Errorf("%w for val '%s'\n", NoMatchErr, val)
}

func (toks Tokens) FindByIndex(ti []int) (Tokens, error) {
	var tokens Tokens
	for _, tok := range ti {
		if tok < toks.Len() {
			tokens = append(tokens, toks[tok])
		}
	}
	if tokens.Len() > 0 {
		return tokens, nil
	}
	return nil, fmt.Errorf("%w for indices %v\n", NoMatchErr, ti)
}

func (toks Tokens) Values() []string {
	vals := make([]string, toks.Len())
	for i, tok := range toks {
		vals[i] = tok.Value
	}
	return vals
}

func (toks Tokens) Labels() []string {
	vals := make([]string, toks.Len())
	for i, tok := range toks {
		vals[i] = tok.Label
	}
	return vals
}

func (toks Tokens) String(i int) string {
	return toks[i].Value
}

func (toks Tokens) Len() int {
	return len(toks)
}

func (toks Tokens) Sort(cmp func(a, b *Token) int, order string) Tokens {
	tokens := toks
	slices.SortStableFunc(tokens, cmp)
	if order == "desc" {
		slices.Reverse(tokens)
	}
	return tokens
}

func (toks Tokens) SortStable(cmp func(a, b *Token) int, order string) Tokens {
	tokens := toks
	slices.SortStableFunc(tokens, cmp)
	if order == "desc" {
		slices.Reverse(tokens)
	}
	return tokens
}

func (toks Tokens) SortAlphaAsc() Tokens {
	return toks.Sort(SortByAlphaFunc, "asc")
}

func (toks Tokens) SortAlphaDesc() Tokens {
	return toks.Sort(SortByAlphaFunc, "desc")
}

func SortByAlphaFunc(a *Token, b *Token) int {
	switch {
	case a.Value > b.Value:
		return 1
	case a.Value == b.Value:
		return 0
	default:
		return -1
	}
}

func newMatch(str string, idx int) fuzzy.Match {
	return fuzzy.Match{
		Str:   str,
		Index: idx,
	}
}
