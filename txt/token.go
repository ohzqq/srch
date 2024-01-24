package txt

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
)

type Tokens struct {
	tokens map[string]*Token
}

func NewTokens() *Tokens {
	return &Tokens{
		tokens: make(map[string]*Token),
	}
}

type Token struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	bits        *roaring.Bitmap
	fuzzy.Match `json:"-"`
}

func NewToken(label string) *Token {
	return &Token{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (f *Token) SetValue(txt string) *Token {
	f.Value = txt
	return f
}

func (f *Token) Count() int {
	return f.bits.GetCardinality()
}
