package srch

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Src struct {
	Data   []any
	fields []string
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

func (r *Src) Search(q string) (*Results, error) {
	res := &Results{}
	if q == "" {
		res.Data = r.Data
		return res, nil
	}

	matches := r.FuzzyFind(q)
	for _, m := range matches {
		res.Data = append(res.Data, r.Data[m.Index])
	}
	return res, nil
}

func (src *Src) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, src)
}

func (r *Src) Len() int {
	return len(r.Data)
}

func (r *Src) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.Data[i]),
		r.fields,
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}
