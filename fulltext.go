package srch

import (
	"strings"
	"unicode"

	"github.com/RoaringBitmap/roaring"
	"github.com/kljensen/snowball/english"
	"github.com/ohzqq/srch/txt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Fulltext struct {
	*Field
}

func FullText(fields []*Field, q string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, field := range fields {
		bits = append(bits, field.Search(q))
	}
	return processBitResults(bits, And)
}

func FulltextAnalyzer(str any) []*txt.Token {
	var tokens []string
	var items []*txt.Token
	for _, token := range strings.FieldsFunc(cast.ToString(str), NotAlphaNumeric) {
		lower := strings.ToLower(token)
		if !lo.Contains(stopWords, lower) {
			items = append(items, txt.NewToken(token))
			tokens = append(tokens, lower)
		}
	}
	for i, t := range stemmerFilter(tokens) {
		items[i].Value = t
	}
	return items
}

func rmStopWords(tokens []string) []string {
	var words []string
	for _, token := range tokens {
		if !lo.Contains(stopWords, token) {
			words = append(words, token)
		}
	}
	return words
}

func stripNonAlphaNumeric(token string) string {
	s := []byte(token)
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}

func lowerCase(tokens []string) []string {
	lower := make([]string, len(tokens))
	for i, str := range tokens {
		lower[i] = strings.ToLower(str)
	}
	return lower
}

func NotAlphaNumeric(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}

func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = english.Stem(token, false)
	}
	return r
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
