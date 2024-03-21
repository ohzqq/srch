package facet

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
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

//func (i *Item) MarshalJSON() ([]byte, error) {
//  //i.Count = i.Len()
//  //i.RelatedTo = i.GetItems()
//  ugh := map[string]any{
//    "count":     i.Len(),
//    "relatedTo": i.GetItems(),
//  }
//  return json.Marshal(ugh)
//}

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

func lowerCase(tokens []string) []string {
	lower := make([]string, len(tokens))
	for i, str := range tokens {
		lower[i] = strings.ToLower(str)
	}
	return lower
}

func stripNonAlphaNumeric(token string) string {
	s := []byte(token)
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}
