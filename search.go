package srch

import (
	"github.com/RoaringBitmap/roaring"
)

func Search(idx *Index, params string) *roaring.Bitmap {
	req := NewQuery(params)

	q := req.Query()
	if q == "" {
		return roaring.New()
	}

	bits := idx.Bitmap()
	bits.And(idx.FuzzySearch(q))
	return bits
}
