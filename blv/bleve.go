package blv

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Index struct {
	*param.SrchCfg
	count int
}

func Open(cfg *param.SrchCfg) *Index {
	blv, err := bleve.Open(cfg.BlvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer blv.Close()

	idx := &Index{
		SrchCfg: cfg,
	}
	idx.count = idx.Count()
	return idx
}

func New(cfg *param.SrchCfg) (*Index, error) {
	idx := &Index{SrchCfg: cfg}
	blv, err := bleve.New(cfg.BlvPath, bleve.NewIndexMapping())
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
	req = bleve.NewSearchRequest(q)

	return idx.search(req)
}

func (idx *Index) search(req *bleve.SearchRequest) ([]map[string]any, error) {
	blv, err := bleve.Open(idx.BlvPath)
	if err != nil {
		return nil, err
	}
	defer blv.Close()

	res, err := blv.Search(req)
	if err != nil {
		return nil, err
	}

	println(res.Total)

	data := make([]map[string]any, res.Hits.Len())
	for i, hit := range res.Hits {
		data[i] = hit.Fields
	}

	//bits := roaring.New()
	//for _, hit := range res.Hits {
	//bits.Add(cast.ToUint32(hit.ID))
	//}
	return data, nil
}

func (idx *Index) Index(uid string, data map[string]any) error {
	blv, err := bleve.Open(idx.BlvPath)
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
	idx.count = idx.Count()

	return nil
}

func (idx *Index) Bitmap() ([]map[string]any, error) {
	q := bleve.NewMatchAllQuery()
	req := bleve.NewSearchRequest(q)
	return idx.search(req)
}

func (idx *Index) Count() int {
	blv, err := bleve.Open(idx.BlvPath)
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
