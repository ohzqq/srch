package blv

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Index struct {
	*param.Params
	count int
}

func Open(params *param.Params) *Index {
	blv, err := bleve.Open(params.Path)
	if err != nil {
		log.Fatalf("%v at %s\n", err, params.Path)
	}
	defer blv.Close()

	idx := &Index{
		Params: params,
	}
	idx.SetCount(blv)
	return idx
}

func New(params *param.Params) (*Index, error) {
	idx := &Index{Params: params}
	blv, err := bleve.New(params.Path, bleve.NewIndexMapping())
	if err != nil {
		return idx, err
	}
	defer blv.Close()

	return idx, nil
}

func (idx *Index) Search(kw string) ([]map[string]any, error) {
	var q query.Query
	q = bleve.NewMatchAllQuery()

	if kw != "" {
		q = bleve.NewTermQuery(kw)
	}

	req := blvReq(q, idx.count)

	blv, err := bleve.Open(idx.Path)
	if err != nil {
		return nil, err
	}
	defer blv.Close()

	req.Fields = idx.SrchAttr
	if idx.Has(param.Facets) {
		req.Fields = append(req.Fields, idx.Facets...)
	}
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

func blvReq(q query.Query, count int) *bleve.SearchRequest {
	return bleve.NewSearchRequestOptions(q, count, 0, true)
}

func search(path string, req *bleve.SearchRequest) ([]map[string]any, error) {
	blv, err := bleve.Open(path)
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

	fmt.Printf("indexing %v\n", uid)

	return blv.Index(uid, data)
}

func (idx *Index) Batch(data []map[string]any) error {
	blv, err := bleve.Open(idx.Path)
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
	idx.SetCount(blv)

	return nil
}

func (idx *Index) SetCount(blv bleve.Index) *Index {
	idx.count = getDocCount(blv)
	return idx
}

func (idx *Index) Len() int {
	return idx.count
}

func getDocCount(blv bleve.Index) int {
	c, err := blv.DocCount()
	if err != nil {
		return 0
	}
	return int(c)
}
