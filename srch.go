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

func FullText(data []map[string]any, q string, fieldNames ...string) []map[string]any {
	fields := IndexText(data, fieldNames, "id")
	return fullTextSearch(data, fields, q)
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

//func (idx *Index) get(q string) (*Index, error) {
//  data := idx.search(q)
//  res := CopyIndex(idx, data)

//  if res.interactive {
//    return res.Choose()
//  }

//  return res.Results()
//}

//func (idx *Index) Results() (*Index, error) {
//  return idx.getResults(), nil
//}

//func (idx *Index) getResults(ids ...int) *Index {
//  if len(ids) > 0 {
//    idx.Data = collectResults(idx.GetData(), ids)
//    return idx
//  }
//  return idx
//}

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

//func (idx *Index) Choose() (*Index, error) {
//  ids, err := Choose(idx)
//  if err != nil {
//    return &Index{}, err
//  }

//  res := idx.getResults(ids...)

//  return res, nil
//}

//func (r *Index) String(i int) string {
//  s := lo.PickByKeys(
//    r.GetData()[i],
//    r.SearchableFields(),
//  )
//  vals := cast.ToStringSlice(lo.Values(s))
//  return strings.Join(vals, "\n")
//}

//func (r *Index) Len() int {
//  return len(r.GetData())
//}

//func (r *Index) FuzzyFind(q string) fuzzy.Matches {
//  return fuzzy.FindFrom(q, r)
//}

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
