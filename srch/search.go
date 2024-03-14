package srch

import (
	"github.com/RoaringBitmap/roaring"
)

type Indexer interface {
	Index(uid string, data ...map[string]any) error
	Searcher
}

type Searcher interface {
	Search(query string) (*roaring.Bitmap, error)
}
