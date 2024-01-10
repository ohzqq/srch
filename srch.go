package srch

import (
	"strings"
	"unicode"

	"github.com/RoaringBitmap/roaring"
	"github.com/kljensen/snowball/english"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type SearchFunc func(string) []map[string]any

func FullText(data []map[string]any, q string, text ...string) *Results {
	idx := New()

	if len(data) < 1 {
		return NewResults(idx, data)
	}

	if len(text) < 1 {
		text = lo.Keys(data[0])
	}

	for _, t := range text {
		idx.AddField(NewTextField(t))
	}

	return idx.Search(q, SliceSrc(data))
}

func FullTextFunc(data []map[string]any, fields []*Field) SearchFunc {
	return func(q string) []map[string]any {
		return fullTextSearch(data, fields, q)
	}
}

func IndexText(data []map[string]any, text []string) []*Field {
	fields := NewTextFields(text)
	return IndexData(data, fields)
}

func fullTextSearch(data []map[string]any, fields []*Field, q string) []map[string]any {
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

//func FuzzySearch(data []map[string]any, fields ...string) SearchFunc {
//  return func(q string) []map[string]any {
//    if q == "" {
//      return data
//    }

//    src := GetSearchableFieldValues(data, fields)
//    var res []map[string]any
//    for _, m := range fuzzy.Find(q, src) {
//      res = append(res, data[m.Index])
//    }
//    return res
//  }
//}

func GetSearchableFieldValues(data []map[string]any, fields []string) []string {
	src := make([]string, len(data))
	for i, d := range data {
		s := lo.PickByKeys(d, fields)
		vals := cast.ToStringSlice(lo.Values(s))
		src[i] = strings.Join(vals, "\n")
	}
	return src
}

func collectResults(d []map[string]any, ids []int) []map[string]any {
	if len(ids) > 0 {
		data := make([]map[string]any, len(ids))
		for i, id := range ids {
			data[i] = d[id]
		}
		return data
	}
	return d
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
