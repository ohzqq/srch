package txt

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Token struct {
	Value string `json:"value"`
	Label string `json:"label"`
	bits  *roaring.Bitmap
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
		bits:  roaring.New(),
	}
}

func (kw *Token) Bitmap() *roaring.Bitmap {
	return kw.bits
}

func (kw *Token) SetValue(txt string) *Token {
	kw.Value = txt
	return kw
}

func (kw *Token) Items() []int {
	i := kw.bits.ToArray()
	return cast.ToIntSlice(i)
}

func (kw *Token) Count() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Len() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Contains(id int) bool {
	return kw.bits.ContainsInt(id)
}

func (kw *Token) Add(ids ...int) {
	for _, id := range ids {
		if !kw.Contains(id) {
			kw.bits.AddInt(id)
		}
	}
}

func (kw *Token) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		"count": kw.Len(),
		"value": kw.Label,
		"hits":  kw.Items(),
	}
	return json.Marshal(item)
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

func newMatch(str string, idx int) fuzzy.Match {
	return fuzzy.Match{
		Str:   str,
		Index: idx,
	}
}
