package srch

import (
	"log"
	"strings"
	"unicode"

	"github.com/RoaringBitmap/roaring"
	"github.com/kljensen/snowball/english"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type SearchFunc func(...any) []map[string]any

func Search(data []map[string]any, fields []*Field, fn SearchFunc, q Query) *Index {
	idx := &Index{
		Data:   data,
		Fields: fields,
		search: fn,
	}
	idx.BuildIndex()
	idx.Data = idx.search(q)
	return idx
}

func FullText(data []map[string]any, fieldNames ...string) SearchFunc {
	return func(query ...any) []map[string]any {
		var q string
		if len(query) > 0 {
			q = cast.ToString(query[0])
		}
		if q == "" {
			return data
		}

		idx := New(SliceSrc(data), WithFields(fieldNames))
		idx.BuildIndex()

		var bits []*roaring.Bitmap
		for _, field := range fieldNames {
			f, err := idx.GetField(field)
			if err != nil {
				log.Fatal(err)
			}
			bits = append(bits, f.Search(q))
		}
		res := roaring.ParOr(viper.GetInt("workers"), bits...)
		return collectResults(idx.Data, cast.ToIntSlice(res))
	}
}

func FuzzySearch(data []map[string]any, fields ...string) SearchFunc {
	return func(qq ...any) []map[string]any {
		var q string
		if len(qq) > 0 {
			q = qq[0].(string)
		}
		if q == "" {
			return data
		}

		src := GetSearchableFieldValues(data, fields)
		var res []map[string]any
		for _, m := range fuzzy.Find(q, src) {
			res = append(res, data[m.Index])
		}
		return res
	}
}

func GetSearchableFieldValues(data []map[string]any, fields []string) []string {
	src := make([]string, len(data))
	for i, d := range data {
		s := lo.PickByKeys(
			cast.ToStringMap(d),
			fields,
		)
		vals := cast.ToStringSlice(lo.Values(s))
		src[i] = strings.Join(vals, "\n")
	}
	return src
}

func (idx *Index) Search(q any) *Index {
	filters, err := NewQuery(q)
	if err != nil {
		log.Fatal(err)
	}
	idx.Query = filters

	res, err := idx.get(filters.Keywords()...)
	if err != nil {
		return idx
	}

	if !res.HasFacets() {
		return res
	}

	res.BuildIndex()

	return Filter(res)
}

func (idx *Index) get(q ...string) (*Index, error) {
	data := idx.search(lo.ToAnySlice(q)...)
	res := CopyIndex(idx, data)

	if res.interactive {
		return res.Choose()
	}

	return res.Results()
}

func (idx *Index) Results() (*Index, error) {
	return idx.getResults(), nil
}

func (idx *Index) getResults(ids ...int) *Index {
	if len(ids) > 0 {
		idx.Data = collectResults(idx.Data, ids)
		return idx
	}
	return idx
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

func (idx *Index) Choose() (*Index, error) {
	ids, err := Choose(idx)
	if err != nil {
		return &Index{}, err
	}

	res := idx.getResults(ids...)

	return res, nil
}

func (r *Index) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.Data[i]),
		r.SearchableFields(),
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}

func (r *Index) Len() int {
	return len(r.Data)
}

func (r *Index) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, r)
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
