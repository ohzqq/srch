package srch

import "github.com/RoaringBitmap/roaring"

type Searcher interface {
	Search(string, ...*roaring.Bitmap) (*Response, error)
}
