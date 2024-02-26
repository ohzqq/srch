package srch

import (
	"github.com/blevesearch/bleve/v2"
)

type FullText struct {
	*bleve.Index
	memOnly bool
	path    string
}

type FTOpt func(*FullText)

func NewTextIndex(opts ...FTOpt) (*bleve.Index, error) {
	ft := &FullText{
		path: "idx",
	}

	for _, opt := range opts {
		opt(ft)
	}

	m := bleve.NewIndexMapping()

	if tf.memOnly {
		return bleve.NewMemOnly(m)
	}

}

func MemOnly(tf *FullText) {
	tf.memOnly = true
}

func FTPath(path string) FTOpt {
	return func(ft FullText) {
		ft.path = path
	}
}
