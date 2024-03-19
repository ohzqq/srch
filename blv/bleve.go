package blv

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Index struct {
	*param.Params
	count int
}

func Open(cfg *param.Params) *Index {
	println(cfg.Path)
	blv, err := bleve.Open(cfg.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer blv.Close()

	idx := &Index{
		Params: cfg,
	}
	c, err := blv.DocCount()
	if err != nil {
		c = 0
	}
	idx.count = int(c)
	return idx
}

func New(cfg *param.Params) (*Index, error) {
	idx := &Index{Params: cfg}
	blv, err := bleve.New(cfg.Path, bleve.NewIndexMapping())
	if err != nil {
		return idx, err
	}
	defer blv.Close()

	return idx, nil
}

func (idx *Index) Search(query string) ([]map[string]any, error) {

	var req *bleve.SearchRequest

	if query == "" {
		q := bleve.NewMatchAllQuery()
		req = bleve.NewSearchRequestOptions(q, idx.count, 0, true)
		return idx.search(req)
	}

	q := bleve.NewTermQuery(query)
	req = bleve.NewSearchRequestOptions(q, idx.count, 0, true)
	return idx.search(req)
}

func (idx *Index) search(req *bleve.SearchRequest) ([]map[string]any, error) {
	blv, err := bleve.Open(idx.Path)
	if err != nil {
		return nil, err
	}
	defer blv.Close()

	req.Fields = []string{"*"}
	res, err := blv.Search(req)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]any, res.Hits.Len())
	for i, hit := range res.Hits {
		data[i] = hit.Fields
	}

	return data, nil
}

func (idx *Index) Index(uid string, data map[string]any) error {
	blv, err := bleve.Open(idx.Path)
	if err != nil {
		return err
	}
	defer blv.Close()

	return blv.Index(uid, data)
}

func (idx *Index) Batch(data []map[string]any) error {
	blv, err := bleve.Open(idx.BlvPath)
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
			if it, ok := doc[idx.UID]; ok {
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
	dc, err := blv.DocCount()
	if err != nil {
		dc = 0
	}
	idx.count = int(dc)

	return nil
}

func (idx *Index) Bitmap() ([]map[string]any, error) {
	q := bleve.NewMatchAllQuery()
	req := bleve.NewSearchRequest(q)
	return idx.search(req)
}

func (idx *Index) Count() int {
	blv, err := bleve.Open(idx.Path)
	if err != nil {
		return 0
	}
	defer blv.Close()

	c, err := blv.DocCount()
	if err != nil {
		return 0
	}

	return int(c)
}

func (idx *Index) Len() int {
	return idx.count
}
