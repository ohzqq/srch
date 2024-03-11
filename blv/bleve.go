package blv

import (
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cast"
)

type Index struct {
	path string
}

func Open(path string) *Index {
	idx := &Index{
		path: path,
	}
	return idx
}

func New(path string) (*Index, error) {
	idx := &Index{
		path: path,
	}

	blv, err := bleve.New(path, bleve.NewIndexMapping())
	if err != nil {
		return idx, err
	}
	defer blv.Close()

	return idx, nil
}

func (idx *Index) Search(query string) (*roaring.Bitmap, error) {
	blv, err := bleve.Open(idx.path)
	if err != nil {
		return nil, err
	}
	defer blv.Close()

	q := bleve.NewTermQuery(query)
	req := bleve.NewSearchRequest(q)
	res, err := blv.Search(req)
	if err != nil {
		return nil, err
	}

	bits := roaring.New()
	for _, hit := range res.Hits {
		bits.Add(cast.ToUint32(hit.ID))
	}
	return bits, nil
}
func (idx *Index) Index(uid string, data ...map[string]any) error {
	switch len(data) {
	case 0:
		return nil
	}

	if len(data) > 1 {
		return idx.Batch(uid, data)
	}

	blv, err := bleve.Open(idx.path)
	if err != nil {
		return err
	}
	defer blv.Close()

	if id, ok := data[0][uid]; ok {
		uid = cast.ToString(id)
	}

	return blv.Index(uid, data[0])
}

func (idx *Index) Batch(uid string, data []map[string]any) error {
	blv, err := bleve.Open(idx.path)
	if err != nil {
		return err
	}
	defer blv.Close()

	batch := blv.NewBatch()

	batchSize := 1000
	total := len(data)
	numB := total / batchSize
	if total%batchSize > 0 {
		numB++
	}
	s := 0
	c := 0
	for b := 1; b < numB+1; b++ {
		e := b * batchSize
		if e > total {
			e = total
		}

		if b < numB+1 {
			fmt.Printf("Indexing batch (%d docs)...\n", e-s)
			err := blv.Batch(batch)
			if err != nil {
				return err
			}
			batch = blv.NewBatch()
		}

		for i := 0; i <= batchSize; {
			if c > total-1 {
				break
			}

			doc := data[c]

			id := cast.ToString(c)
			if it, ok := doc[uid]; ok {
				id = cast.ToString(it)
			}

			err = batch.Index(id, doc)
			if err != nil {
				return err
			}

			c++
		}

		s += batchSize
		e += batchSize
	}

	return nil
}

func (idx *Index) Len() int {
	blv, err := bleve.Open(idx.path)
	if err != nil {
		return 0
	}
	defer blv.Close()

	c, _ := blv.DocCount()
	return int(c)
}
