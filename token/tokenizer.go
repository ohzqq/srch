package token

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

type Analyzer int

const (
	Keywords Analyzer = iota
	Fulltext
)

func (t Analyzer) Tokenize(og ...string) []string {
	switch t {
	case Fulltext:
		return TokenizeFulltext(og)
	case Keywords:
		return TokenizeKeywords(og)
	default:
		return og
	}
}

func TokenizeKeywords(og []string) []string {
	toks := make([]string, len(og))
	for i, t := range og {
		toks[i] = strings.ToLower(t)
	}
	return toks
}

func TokenizeFulltext(og []string) []string {
	var toks []string
	for _, v := range og {
		tokens := Split(v)
		tokens = RemoveStopwords(tokens...)
		for _, t := range tokens {
			t = Stem(t)
			toks = append(toks, Normalize(t))
		}
	}
	return toks
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
