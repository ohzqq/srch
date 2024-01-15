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
	idx := New(q)

	if len(data) < 1 {
		return idx
	}

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
	return FilterDataByID(data, cast.ToIntSlice(res.ToArray()))
}

func Tokenizer(str string) []*FacetItem {
	var tokens []string
	var items []*FacetItem
	for _, token := range strings.FieldsFunc(str, NotAlphaNumeric) {
		lower := strings.ToLower(token)
		if !lo.Contains(stopWords, lower) {
			items = append(items, NewFacetItem(token))
			tokens = append(tokens, lower)
		}
	}
	for i, t := range stemmerFilter(tokens) {
		items[i].Value = t
	}
	return items
}

func FacetTokenizer(val any) []*FacetItem {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*FacetItem, len(tokens))
	for i, token := range tokens {
		items[i] = NewFacetItem(token)
		items[i].Value = normalizeText(token)
	}
	return items
}

func normalizeText(token string) string {
	fields := lowerCase(strings.Split(token, " "))
	for t, term := range fields {
		if len(term) == 1 {
			fields[t] = term
		} else {
			fields[t] = stripNonAlphaNumeric(term)
		}
	}
	return strings.Join(fields, " ")
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
