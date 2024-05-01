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
	var tokens []string
	for _, str := range og {
		tokens = append(tokens, t.analyze(str)...)
	}
	//slices.Sort(tokens)
	return lo.Uniq(tokens)
}

func (t Analyzer) analyze(str string) []string {
	str = strings.ToLower(str)

	tokens := t.split(str)
	tokens = t.rmStopwords(tokens)
	tokens = t.stem(tokens)
	return tokens
}

func (t Analyzer) rmStopwords(toks []string) []string {
	toks = RemovePunct(toks)
	switch t {
	case Simple:
		return RemoveStopLetters(toks)
	case Standard:
		return RemoveStopwords(toks)
	default:
		return toks
	}
}

func (t Analyzer) split(str string) []string {
	switch t {
	case Keyword:
		return []string{str}
	default:
		return SplitOnWhitespaceAndPunct(str)
	}
}

func (t Analyzer) stem(tokens []string) []string {
	switch t {
	case Keyword:
		return tokens
	default:
		return StemWords(tokens)
	}
}

func ToLower(toks []string) []string {
	low := make([]string, len(toks))
	for i, tok := range toks {
		low[i] = strings.ToLower(tok)
	}
	return low
}

func RemovePunct(toks []string) []string {
	var none []string
	for _, tok := range toks {
		if len(tok) > 1 {
			none = append(none, tok)
		} else {
			if r := rune(tok[0]); !unicode.IsPunct(r) {
				none = append(none, tok)
			}
		}
	}
	return none
}

func SplitAndNormalize(tok string) []string {
	var toks []string
	for _, t := range SplitOnWhitespaceAndPunct(tok) {
		toks = append(toks, t)
	}
	return RemoveStopLetters(toks)
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

func RemoveStopWords(tokens ...string) []string {
	return lo.Without(tokens, defaultStopwords...)
}

func RemoveStopwords(tokens []string) []string {
	return lo.Without(tokens, defaultStopwords...)
}

func RemoveStopLetters(tokens []string) []string {
	return lo.Without(tokens, stopLetters...)
}

func StemWords(toks []string) []string {
	stem := make([]string, len(toks))
	for i, tok := range toks {
		stem[i] = Stem(tok)
	}
	return stem
}

func Stem(tok string) string {
	return english.Stem(tok, false)
}

var stopLetters = []string{
	"t",
	"s",
	"ll",
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
