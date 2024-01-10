package srch

import (
	"strings"
	"unicode"

	"github.com/RoaringBitmap/roaring"
	"github.com/kljensen/snowball/english"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func FullText(data []map[string]any, q string, fields ...string) *Index {
	idx := New()

	if len(data) < 1 {
		return idx
	}

	idx.AddField(GetFieldsFromSlice(data, fields)...)
	idx.Index(data)

	return idx.Search(q)
}

func FullTextSrchFunc(data []map[string]any, fields []*Field) SearchFunc {
	return func(q string) []map[string]any {
		return searchFullText(data, fields, q)
	}
}

func searchFullText(data []map[string]any, fields []*Field, q string) []map[string]any {
	if q == "" {
		return data
	}
	var bits []*roaring.Bitmap
	for _, field := range fields {
		bits = append(bits, field.Search(q))
	}
	res := processBitResults(bits, "and")
	return collectResults(data, cast.ToIntSlice(res.ToArray()))
}

func Tokenizer(str string) []string {
	var tokens []string
	for _, token := range strings.FieldsFunc(str, splitFields) {
		token := strings.ToLower(token)
		if !lo.Contains(stopWords, token) {
			tokens = append(tokens, token)
		}
	}
	return stemmerFilter(tokens)
}

func splitFields(c rune) bool {
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
