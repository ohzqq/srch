package srch

import (
	"github.com/RoaringBitmap/roaring"
)

type Indexer interface {
	New(settings string) (Searcher, error)
	Open(settings string) (Searcher, error)
	Index(uid string, data ...map[string]any) error
}

type Searcher interface {
	Search(query string) (*roaring.Bitmap, error)
}
