package blv

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cast"
)

type Index struct {
	bleve.Index
	memOnly bool
	path    string
}

type FTOpt func(*Index)

func New(path string) (*Index, error) {
	idx := &Index{
		path: path,
	}

	b, err := bleve.New(path, bleve.NewIndexMapping())
	if err != nil {
		return idx, err
	}

	idx.Index = b
	return idx
}

func (idx *Index) Open(path string) *Index {
}

func (idx *Index) Index(uid string, data ...map[string]any) error {
	batchSize := 1000
	i := 0
	batch := idx.Index.NewBatch()
	for di, b := range data {
		if i%batchSize == 0 {
			fmt.Printf("Indexing batch (%d docs)...\n", i)
			err := idx.Index.Batch(batch)
			if err != nil {
				return err
			}
			batch = idx.Index.NewBatch()
		}

		id := cast.ToString(di)
		if it, ok := b[uid]; ok {
			id = cast.ToString(it)
		}

		err = batch.Index(id, b)
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func SearchBleve(path, query string) (*bleve.SearchResult, error) {
	blv, err := bleve.Open(path)
	if err != nil {
		return nil, err
	}
	defer blv.Close()

	//q := bleve.NewQueryStringQuery(query)
	q := bleve.NewTermQuery(query)
	req := bleve.NewSearchRequest(q)
	return blv.Search(req)
}

func NewMemOnly(fd string) (bleve.Index, error) {
	idx, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return nil, err
	}

	err = BatchIndex(idx, fd)
	if err != nil {
		return nil, err
	}

	return idx, nil
}

func MemOnly(tf *Index) {
	tf.memOnly = true
}

func FTPath(path string) FTOpt {
	return func(ft *Index) {
		ft.path = path
	}
}
