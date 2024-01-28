package srch

import (
	"encoding/json"
	"net/url"
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/txt"
	"github.com/spf13/viper"
)

const (
	Or          = "or"
	And         = "and"
	Not         = `not`
	FacetField  = `facet`
	SortByCount = `count`
	SortByAlpha = `alpha`
)

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	SortBy    string
	Order     string
	*txt.Tokens
}

func NewField(attr string) *Field {
	f := &Field{
		Sep:    ".",
		SortBy: "count",
		Order:  "desc",
		Tokens: txt.NewTokens(),
	}
	parseAttr(f, attr)
	return f
}

func (f *Field) MarshalJSON() ([]byte, error) {
	field := map[string]any{
		"attribute": f.Attribute,
		"sort_by":   f.SortBy,
		"order":     f.Order,
		"items":     f.GetTokens(),
	}

	d, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (f *Field) GetTokens() []*txt.Token {
	return f.Tokens.Tokens()
}

func (f *Field) Find(kw string) []*txt.Token {
	return f.Tokens.Find(kw)
}

func (f *Field) And(filters []string) *roaring.Bitmap {
	//var bits []*roaring.Bitmap
	var bits *roaring.Bitmap
	for _, filt := range filters {
		var not string
		var ok bool
		var tokens []*txt.Token
		not, ok = IsNegative(filt)
		if ok {
			tokens = f.Find(not)
		} else {
			tokens = f.Find(filt)
		}
		for i, token := range tokens {
			if i == 0 {
				bits = token.Bitmap()
			}
			switch ok {
			case true:
				bits.AndNot(token.Bitmap())
			default:
				bits.And(token.Bitmap())
			}
			//bits = append(bits, token.Bitmap())
		}
	}
	return bits
	//return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (f *Field) Filterz(bits *roaring.Bitmap, vals url.Values) *roaring.Bitmap {
	//var not []string
	for op, filters := range vals {
		for _, filt := range filters {
			var ok bool
			filt, ok = IsNegative(filt)
			if !ok {
				switch op {
				case Or:
					bits.Or(f.Tokens.Filter(filt))
				case And:
					bits.And(f.Tokens.Filter(filt))
				}

			}
		}
	}
	return bits
}

func (f *Field) Or(filters ...string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, filt := range filters {
		for _, token := range f.Find(filt) {
			//bits.Or(token.Bitmap())
			bits = append(bits, token.Bitmap())
		}
	}
	//return bits
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func parseAttr(field *Field, attr string) {
	i := 0
	for attr != "" {
		var a string
		a, attr, _ = strings.Cut(attr, ":")
		if a == "" {
			continue
		}
		switch i {
		case 0:
			field.Attribute = a
		case 1:
			field.SortBy = a
		case 2:
			field.Order = a
		}
		i++
	}
}

func SortItemsByCount(items []*txt.Token) []*txt.Token {
	slices.SortStableFunc(items, SortByCountFunc)
	return items
}

func SortItemsByLabel(items []*txt.Token) []*txt.Token {
	slices.SortStableFunc(items, SortByLabelFunc)
	return items
}

func SortByCountFunc(a *txt.Token, b *txt.Token) int {
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

func SortByLabelFunc(a *txt.Token, b *txt.Token) int {
	switch {
	case a.Label > b.Label:
		return 1
	case a.Label == b.Label:
		return 0
	default:
		return -1
	}
}
