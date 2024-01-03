package srch

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Source interface {
	Get(...string) ([]any, error)
}

type Src struct {
	Data    []any `json:"data"`
	fields  []string
	matches fuzzy.Matches
}

func NewSrc(data []any, fields ...string) *Src {
	f := []string{"title"}
	if len(fields) > 0 {
		f = fields
	}
	return &Src{
		Data:   data,
		fields: f,
	}
}

func (r *Src) Search(q string) error {
	if q == "" {
		return NewResults(r.Data), nil
	}

	var data []any
	r.matches = r.FuzzyFind(q)
	for _, m := range r.matches {
		data = append(data, r.Data[m.Index])
	}
	return NewResults(data), nil
}

func (src *Src) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, src)
}

func (src *Src) Find(q string) Results {
	return src
}

func (r *Src) Len() int {
	return len(r.Data)
}

func (r *Src) Matches() []int {
	if len(r.matches) > 0 {
		fn := func(m fuzzy.Match, _ int) int {
			return m.Index
		}
		return lo.Map(r.matches, fn)
	}
	return r.dataIDs()
}

func (r *Src) dataIDs() []int {
	fn := func(_ any, index int) int {
		return index
	}
	return lo.Map(r.Data, fn)
}

func (r *Src) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.Data[i]),
		r.fields,
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}
