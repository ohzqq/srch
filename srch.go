package srch

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type SearchFunc func(string) []any

func FuzzySearch(data []any, fields ...string) SearchFunc {
	return func(q string) []any {
		if q == "" {
			return data
		}

		src := GetSearchableFieldValues(data, fields)
		var res []any
		for _, m := range fuzzy.Find(q, src) {
			res = append(res, data[m.Index])
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
	res := CopyIndex(s, s.search(q))

	if res.interactive {
		return res.Choose()
	}

	return res.Results()
}

func (m *Index) Results() (*Index, error) {
	return m.getResults(), nil
}

func (m *Index) getResults(ids ...int) *Index {
	if len(ids) > 0 {
		m.Data = make([]any, len(ids))
		for i, id := range ids {
			m.Data[i] = m.Data[id]
		}
		return m
	}

	return m
}

func (s *Index) Choose() (*Index, error) {
	ids, err := Choose(s)
	if err != nil {
		return &Index{}, err
	}

	res := s.getResults(ids...)

	return res, nil
}

func (r *Index) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.Data[i]),
		r.SearchableFields,
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
