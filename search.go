package srch

import (
	"github.com/RoaringBitmap/roaring"
)

func Search(idx *Index, kw string) *roaring.Bitmap {
	return idx.res
}
