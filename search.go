package srch

import (
	"github.com/RoaringBitmap/roaring"
)

type Indexer interface {
	New(settings string) (*Index, error)
	Open(settings string) (*Index, error)
	Index(uid string, data ...map[string]any) error
}

type Searcher interface {
	Search(query string) (*roaring.Bitmap, error)
}
