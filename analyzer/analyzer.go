package analyzer

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
	"github.com/samber/lo"
)

type Analyzer int

const (
	Keyword Analyzer = iota
	Simple
	Standard
)

func (t Analyzer) Tokenize(og ...string) []string {
	switch t {
	case Standard:
		return TokenizeFulltext(og)
	case Keyword:
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
		tokens := SplitAndNormalize(v)
		tokens = RemoveStopWords(tokens...)
		for _, t := range tokens {
			t = Stem(t)
			toks = append(toks, t)
		}
	}
	return toks
}

func TokenizeSimple(og []string) []string {
	var toks []string
	for _, v := range og {
		tokens := SplitAndNormalize(v)
		for _, t := range tokens {
			t = Stem(t)
			toks = append(toks, t)
		}
	}
	return toks
}

func Normalize(t string) string {
	t = AlphaNumericOnly(t)
	t = strings.ToLower(t)
	return t
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

func SplitOnWhitespaceAndPunct(tok string) []string {
	fn := func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	}
	return strings.FieldsFunc(tok, fn)
}

func SplitAndNormalize(tok string) []string {
	var toks []string
	for _, t := range SplitOnWhitespaceAndPunct(tok) {
		toks = append(toks, t)
	}
	return RemoveStopLetters(toks...)
}

func RemoveStopWords(tokens ...string) []string {
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
	"did",
	"do",
	"does",
	"doing",
	"don",
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
