package srch

import (
	"github.com/RoaringBitmap/roaring"
)

func Search(idx *Index, kw string) *roaring.Bitmap {
	bits := idx.Bitmap()

	if kw == "" {
		return bits
	}

	idx.res.And(idx.FuzzySearch(kw))
	bits.And(idx.FuzzySearch(kw))
	return bits
}
