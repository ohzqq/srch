package srch

import (
	"github.com/RoaringBitmap/roaring"
)

type Indexer interface {
	Open(settings string) (Searcher, error)
	Index(uid string, data ...map[string]any) error
	Searcher
}

type Searcher interface {
	Search(query string) (*roaring.Bitmap, error)
}
