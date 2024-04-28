package facet

import (
	"strings"
	"unicode"

	"github.com/RoaringBitmap/roaring"
	"github.com/kljensen/snowball/english"
	"github.com/spf13/cast"
)

type Item struct {
	Value     string          `json:"value"`
	Label     string          `json:"label"`
	Count     int             `json:"count"`
	RelatedTo []int           `json:"relatedTo"`
	Children  *Facet          `json:"-"`
	bits      *roaring.Bitmap `json:"-"`
}

func NewItem(label string) *Item {
	return &Item{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (i *Item) Bitmap() *roaring.Bitmap {
	return i.bits
}

func (i *Item) SetValue(txt string) *Item {
	i.Value = txt
	return i
}

func (i *Item) GetItems() []int {
	items := i.bits.ToArray()
	return cast.ToIntSlice(items)
}

func (i *Item) Len() int {
	return int(i.bits.GetCardinality())
}

func (i *Item) Contains(id int) bool {
	return i.bits.ContainsInt(id)
}

func (i *Item) Add(ids ...int) {
	for _, id := range ids {
		if !i.Contains(id) {
			i.bits.AddInt(id)
		}
	}
}

func KeywordTokenizer(val any) []*Item {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*Item, len(tokens))
	for i, token := range tokens {
		items[i] = NewItem(token)
		items[i].Value = normalizeText(token)
	}
	return items
}

func Tokenize(vals ...string) []string {
	var toks []string
	for _, v := range vals {
		//toks = append(toks, v)
		tokens := split(v)
		tokens = removeStopwords(tokens)
		for _, t := range tokens {
			toks = append(toks, normalizeStr(t))
		}
	}
	return toks
}

func Keywords(vals ...string) []string {
	var toks []string
	for _, v := range vals {
		toks = append(toks, normalizeKeyword(v))
	}
	return toks
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

func split(tok string) []string {
	fn := func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	}
	return strings.FieldsFunc(tok, fn)
}

func normalizeKeyword(tok string) string {
	return strings.ToLower(tok)
}

func normalizeStr(tok string) string {
	tok = stripNonAlphaNumeric(tok)
	tok = stem(tok)
	return tok
}

func stem(tok string) string {
	return english.Stem(tok, false)
}

func lowerCase(tokens []string) []string {
	lower := make([]string, len(tokens))
	for i, str := range tokens {
		lower[i] = strings.ToLower(str)
	}
	return lower
}

func removeStopwords(tokens []string) []string {
	var toks []string
	for _, t := range tokens {
		if len(t) > 2 {
			toks = append(toks, t)
		}
	}
	return toks
}

func stripNonAlphaNumeric(token string) string {
	s := []byte(token)
	n := 0
	for _, b := range s {
		r := rune(b)
		if unicode.IsLetter(r) ||
			unicode.IsNumber(r) ||
			unicode.IsSpace(r) {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}
