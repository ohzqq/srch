package srch

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/txt"
)

type Token struct {
	*txt.Token
	bits *roaring.Bitmap
}

func NewToken(label, val string) *Token {
	return &Token{
		Token: txt.NewToken(label, val),
		bits:  roaring.New(),
	}
}

func newTok(t *txt.Token) *Token {
	return &Token{
		Token: t,
		bits:  roaring.New(),
	}
}

func (f *Token) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Token) SetValue(txt string) *Token {
	f.Value = txt
	return f
}

func (f *Token) Count() int {
	return int(f.bits.GetCardinality())
}

func (f *Token) Contains(id int) bool {
	return f.bits.ContainsInt(id)
}

func (f *Token) Add(ids ...int) {
	for _, id := range ids {
		if !f.Contains(id) {
			f.bits.AddInt(id)
		}
	}
}

func (f *Token) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		f.Label: f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}
