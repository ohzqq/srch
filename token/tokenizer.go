package token

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

type Tokenizer interface {
	Tokenize(...string) []string
}

func Normalize(tok string) string {
	var toks []string
	for _, t := range Split(tok) {
		t = AlphaNumericOnly(t)
		t = strings.ToLower(t)
		toks = append(toks, t)
	}
	return strings.Join(toks, " ")
}

func AlphaNumericOnly(token string) string {
	s := []byte(token)
	n := 0
	for _, b := range s {
		r := rune(b)
		if unicode.IsLetter(r) ||
			unicode.IsSpace(r) ||
			unicode.IsNumber(r) {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}

func SplitOnWhitespace(tok string) []string {
	return strings.FieldsFunc(tok, unicode.IsSpace)
}

func Split(tok string) []string {
	fn := func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	}
	return strings.FieldsFunc(tok, fn)
}

func RemoveStopwords(tokens ...string) []string {
	var toks []string
	for _, t := range tokens {
		if len(t) > 2 {
			toks = append(toks, t)
		}
	}
	return toks
}

func Stem(tok string) string {
	return english.Stem(tok, false)
}

func normalizeStr(tok string) string {
	tok = strings.ToLower(tok)
	tok = AlphaNumericOnly(tok)
	return tok
}
