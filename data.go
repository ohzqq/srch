package srch

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Data interface {
	fuzzy.Source
	Data() []any
}

type Src struct {
	data   []any
	fields []string
}

func NewSrc() *Src {
	return &Src{
		fields: []string{"title"},
	}
}

func (r *Src) Search(data []any, q string) (*Results, error) {
	r.data = data
	res := &Results{}
	if q == "" {
		res.Data = data
		return res, nil
	}
	matches := r.FuzzyFind(q)
	for _, m := range matches {
		res.Data = append(res.Data, data[m.Index])
		//res.Data = append(res.Data, &FacetItem{Match: m})
	}
	return res, nil
}

func (src *Src) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, src)
}

func (r *Src) Data() []any {
	return r.data
}

func (r *Src) Len() int {
	return len(r.data)
}

func (r *Src) String(i int) string {
	s := lo.PickByKeys(
		cast.ToStringMap(r.data[i]),
		r.fields,
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}
