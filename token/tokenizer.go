package token

import (
	"strings"
	"unicode"
)

type Tokenizer interface {
	Tokenize(string) []string
}

func Normalize(tok string) string {
}

func AlphaNumericOnly(token string) string {
	s := []byte(token)
	n := 0
	for _, b := range s {
		r := rune(b)
		if unicode.IsLetter(r) ||
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
