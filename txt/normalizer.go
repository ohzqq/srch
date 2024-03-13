package txt

import (
	"strings"

	"github.com/kljensen/snowball/english"
	"github.com/ohzqq/txt/sep"
)

type Normalizer func(string) string

func Split(str string, s sep.Func) []string {
	if s == nil {
		return []string{str}
	}
	return strings.FieldsFunc(str, s)
}

func Normalize(ls []string, normies []Normalizer) Tokens {
	toks := make(Tokens, len(ls))
	for i, l := range ls {
		t := normalize(l, normies)
		toks[i] = NewToken(l, t)
	}
	return toks
}

func normalizeSlice(ls []string, normies []Normalizer) []string {
	for i, l := range ls {
		ls[i] = normalize(l, normies)
	}
	return ls
}

func normalize(label string, normies []Normalizer) string {
	for _, norm := range normies {
		label = norm(label)
	}
	return label
}

func StripPunct(token string) string {
	var s []byte
	for _, b := range []byte(token) {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			s = append(s, b)
		}
	}
	return string(s)
}

func Stem(token string) string {
	return english.Stem(token, false)
}

func DefaultStopWords() []string {
	return stopWords
}

var stopWords = []string{
	"i",
	"vol",
	"what",
	"which",
	"who",
	"whom",
	"this",
	"that",
	"am",
	"is",
	"are",
	"was",
	"were",
	"be",
	"been",
	"being",
	"have",
	"has",
	"had",
	"having",
	"do",
	"does",
	"did",
	"doing",
	"a",
	"an",
	"the",
	"and",
	"but",
	"if",
	"or",
	"because",
	"as",
	"of",
	"at",
	"by",
	"for",
	"with",
	"into",
	"to",
	"from",
	"then",
	"when",
	"where",
	"why",
	"how",
	"no",
	"not",
	"than",
	"too",
}
