package srch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cast"
)

//func TestDocMapping(t *testing.T) {
//  doc := DocMapping(testFieldMappings())
//  d, err := json.Marshal(doc)
//  if err != nil {
//    t.Error(err)
//  }
//  println(string(d))
//}

const blevePath = `testdata/poot`

func TestNewBleveIndex(t *testing.T) {
	idx, err := NewTextIndex(FTPath(blevePath))
	//idx, err := NewTextIndex(MemOnly)
	if err != nil {
		//t.Fatal(err)
		t.Skipf("%v\n", err)
	}

	d, err := idx.DocCount()
	if err != nil {
		t.Error(err)
	}
	println(d)
}

func TestBatchIndex(t *testing.T) {
	batchSize := 1000

	idx, err := bleve.Open(blevePath)
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()

	file, err := os.Open("testdata/audiobks.json")
	if err != nil {
		t.Error(err)
	}

	i := 0
	batch := idx.NewBatch()

	r := bufio.NewReader(file)

	for {
		if i%batchSize == 0 {
			fmt.Printf("Indexing batch (%d docs)...\n", i)
			err := idx.Batch(batch)
			if err != nil {
				t.Error(err)
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
			t.Errorf("error parsing JSON: %v", err)
		}

		book := cast.ToStringMap(doc)

		docID := cast.ToString(book["id"])
		//docID := cast.ToString(i)
		err = batch.Index(docID, book)
		if err != nil {
			t.Error(err)
		}
		i++
	}

	err = file.Close()
	if err != nil {
		t.Error(err)
	}

}

func cleanIdx() {
	idxMeta := filepath.Join(blevePath, "index_meta.json")
	err := os.Remove(idxMeta)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	idxStore := filepath.Join(blevePath, "store")
	err = os.RemoveAll(idxStore)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
