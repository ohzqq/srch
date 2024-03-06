package srch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cast"
)

type FullText struct {
	bleve.Index
	memOnly bool
	path    string
}

type FTOpt func(*FullText)

func NewTextIndex(opts ...FTOpt) (*FullText, error) {
	ft := &FullText{
		path: "idx",
	}

	for _, opt := range opts {
		opt(ft)
	}

	m := bleve.NewIndexMapping()

	var idx bleve.Index
	var err error
	if ft.memOnly {
		idx, err = bleve.NewMemOnly(m)
		if err != nil {
			return nil, err
		}
	}

	//idx, err = bleve.New(ft.path, m)
	idx, err = bleve.Open(ft.path)
	if err != nil {
		return nil, err
	}

	ft.Index = idx
	return ft, nil
}

func BatchIndex(idx bleve.Index, fd string) error {
	batchSize := 1000
	file, err := os.Open("testdata/ndbooks.json")
	if err != nil {
		return err
	}

	i := 0
	batch := idx.NewBatch()

	r := bufio.NewReader(file)

	for {
		if i%batchSize == 0 {
			fmt.Printf("Indexing batch (%d docs)...\n", i)
			err := idx.Batch(batch)
			if err != nil {
				return err
			}
			batch = idx.NewBatch()
		}

		b, _ := r.ReadBytes('\n')
		if len(b) == 0 {
			break
		}

		var doc interface{}
		doc = b
		var err error
		err = json.Unmarshal(b, &doc)
		if err != nil {
			return fmt.Errorf("error parsing JSON: %v", err)
		}

		book := cast.ToStringMap(doc)

		//docID := cast.ToString(book["id"])
		docID := cast.ToString(i)
		err = batch.Index(docID, book)
		if err != nil {
			return err
		}
		i++
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func MemOnly(tf *FullText) {
	tf.memOnly = true
}

func FTPath(path string) FTOpt {
	return func(ft *FullText) {
		ft.path = path
	}
}
