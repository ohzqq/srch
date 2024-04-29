package token

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
	"github.com/samber/lo"
)

type Analyzer int

const (
	Keywords Analyzer = iota
	Simple
	Fulltext
)

func (t Analyzer) Tokenize(og ...string) []string {
	switch t {
	case Fulltext:
		return TokenizeFulltext(og)
	case Keywords:
		return TokenizeKeywords(og)
	case Simple:
		return TokenizeSimple(og)
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

func TokenizeSimple(og []string) []string {
	var toks []string
	for _, v := range og {
		tokens := Split(v)
		tokens = RemoveStopLetters(tokens...)
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
	return lo.Without(tokens, defaultStopwords...)
}

func RemoveStopLetters(tokens ...string) []string {
	return lo.Without(tokens, stopLetters...)
}

func Stem(tok string) string {
	return english.Stem(tok, false)
}

var stopLetters = []string{
	"t",
	"s",
}

var defaultStopwords = []string{
	"a",
	"about",
	"above",
	"after",
	"again",
	"against",
	"all",
	"am",
	"an",
	"and",
	"any",
	"are",
	"as",
	"at",
	"be",
	"because",
	"been",
	"before",
	"being",
	"below",
	"between",
	"both",
	"but",
	"by",
	"can",
	"cant",
	"did",
	"do",
	"does",
	"doing",
	"don",
	"dont",
	"down",
	"during",
	"each",
	"few",
	"for",
	"from",
	"further",
	"had",
	"has",
	"have",
	"having",
	"he",
	"her",
	"here",
	"hers",
	"herself",
	"him",
	"himself",
	"his",
	"how",
	"i",
	"if",
	"in",
	"into",
	"is",
	"it",
	"its",
	"itself",
	"just",
	"me",
	"more",
	"most",
	"my",
	"myself",
	"no",
	"nor",
	"not",
	"now",
	"of",
	"off",
	"on",
	"once",
	"only",
	"or",
	"other",
	"our",
	"ours",
	"ourselves",
	"out",
	"over",
	"own",
	"s",
	"same",
	"she",
	"should",
	"so",
	"some",
	"such",
	"t",
	"than",
	"that",
	"the",
	"their",
	"theirs",
	"them",
	"themselves",
	"then",
	"there",
	"these",
	"they",
	"this",
	"those",
	"through",
	"to",
	"too",
	"under",
	"until",
	"up",
	"very",
	"was",
	"we",
	"were",
	"what",
	"when",
	"where",
	"which",
	"while",
	"who",
	"whom",
	"why",
	"will",
	"with",
	"you",
	"your",
	"yours",
	"yourself",
	"yourselves",
}
