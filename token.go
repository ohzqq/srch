package srch

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/txt"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cast"
)

type Keyword struct {
	*Field
}

func KeywordAnalyzer(val any) []*txt.Token {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*txt.Token, len(tokens))
	for i, token := range tokens {
		items[i] = txt.NewToken(token)
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

func NewFacet(attr string, params ...*Params) *Field {
	f := NewField(attr, params...)
	f.FieldType = FacetField
	return f
}

// Token is a data structure for a Facet's item.
type Token struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	bits        *roaring.Bitmap
	fuzzy.Match `json:"-"`
}

func NewToken(label string) *Token {
	return &Token{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (f *Token) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		"value": f.Value,
		"label": f.Label,
		"count": f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (f *Token) Count() int {
	return len(f.bits.ToArray())
}

func SortItemsByCount(items []*txt.Token) []*txt.Token {
	slices.SortStableFunc(items, sortByCountFunc)
	return items
}

func SortItemsByLabel(items []*txt.Token) []*txt.Token {
	slices.SortStableFunc(items, sortByLabelFunc)
	return items
}

func sortByCountFunc(a *txt.Token, b *txt.Token) int {
	aC := a.Count()
	bC := b.Count()
	switch {
	case aC < bC:
		return 1
	case aC == bC:
		return 0
	default:
		return -1
	}
}

func sortByLabelFunc(a *txt.Token, b *txt.Token) int {
	switch {
	case a.Label < b.Label:
		return 1
	case a.Label == b.Label:
		return 0
	default:
		return -1
	}
}
