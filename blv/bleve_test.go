package blv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

//func TestDocMapping(t *testing.T) {
//  doc := DocMapping(testFieldMappings())
//  d, err := json.Marshal(doc)
//  if err != nil {
//    t.Error(err)
//  }
//  println(string(d))
//}

const blevePath = `../testdata/poot`
const testDataFile = `../testdata/ndbooks.json`

func TestNewBleveIndex(t *testing.T) {
	t.SkipNow()
	_, err := New(blevePath)
	if err != nil {
		t.Error(err)
	}
}

func TestBatchIndex(t *testing.T) {
	//t.SkipNow()
	books := loadData(t)
	println(len(books))

	idx := Open(blevePath)
	err := idx.Index("id", books...)
	if err != nil {
		t.Error(err)
	}
}

func TestOpenIndex(t *testing.T) {
	idx := Open(blevePath)
	docs := idx.Len()
	println(docs)
	if docs != 7000 {
		t.Errorf("got %d docs, expected %d\n", docs, 7252)
	}
}

func TestBleveSearch(t *testing.T) {
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

func loadData(t *testing.T) []map[string]any {
	d, err := os.Open(testDataFile)
	if err != nil {
		t.Error(err)
	}
	defer d.Close()

	var books []map[string]any

	scanner := bufio.NewScanner(d)
	for scanner.Scan() {
		b := make(map[string]any)
		err = json.Unmarshal(scanner.Bytes(), &b)
		if err != nil {
			t.Error(err)
		}
		books = append(books, b)
	}

	return books
}
