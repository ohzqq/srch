package srch

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Searcher interface {
	Search(string) ([]Item, error)
}

type SearchFunc func(string) []any

func FuzzyFind(data []any, fields ...string) SearchFunc {
	return func(q string) []any {
		if q == "" {
			return data
		}

		src := GetSearchableFieldValues(data, fields)
		var res []any
		for _, m := range fuzzy.Find(q, src) {
			res = append(res, m.Index)
		}
		return res
	}
}

func GetSearchableFieldValues(data []any, fields []string) []string {
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

func (s *Index) get(q string) (*Index, error) {
	s.results = s.search(q)

	if s.interactive {
		return s.Choose()
	}

	return s.Results()
}

func (m *Index) Results() (*Index, error) {
	return m.getResults(), nil
}

func (m *Index) getResults(ids ...int) *Index {
	r := &Index{
		//Query: m.query,
	}

	if len(ids) > 0 {
		r.Data = make([]any, len(ids))
		for i, id := range ids {
			r.Data[i] = m.results[id]
		}
		return r
	}
	r.Data = lo.ToAnySlice(m.results)

	return r
}

func (s *Index) Choose() (*Index, error) {
	ids, err := Choose(s.results)
	if err != nil {
		return &Index{}, err
	}

	res := s.getResults(ids...)

	return res, nil
}
