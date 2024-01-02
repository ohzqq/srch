package srch

import (
	"fmt"
	"net/url"

	"github.com/sahilm/fuzzy"
)

type Results struct {
	Data    []any      `json:"data"`
	Facets  []*Facet   `json:"facets"`
	query   Query      `json:"query"`
	Filters url.Values `json:"filters"`
}

type Item interface {
	fmt.Stringer
}

func (r Results) Search(q Query) Results {
	return r
}

func (r Results) String(i int) string {
	return fmt.Sprintf("%v", r.Data[i])
}

func (r Results) Len() int {
	return len(r.Data)
}

func (r Results) FuzzyFind(q string) fuzzy.Matches {
	return fuzzy.FindFrom(q, r)
}
