package blv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

const blevePath = `../testdata/poot.bleve`
const testDataFile = `../testdata/nddata/ndbooks.ndjson`
const cfgStr = `searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`

func TestNewBleveIndex(t *testing.T) {
	//t.SkipNow()
	cleanIdx()

	params, err := param.Parse(cfgStr)
	if err != nil {
		t.Error(err)
	}

	_, err = New(params.SrchCfg)
	if err != nil {
		t.Error(err)
	}
}

func TestBatchIndex(t *testing.T) {
	//t.SkipNow()
	books := loadData(t)

	params, err := param.Parse(cfgStr)
	if err != nil {
		t.Error(err)
	}

	idx := Open(params.SrchCfg)
	err = idx.Batch(books)
	if err != nil {
		t.Error(err)
	}

	if idx.Len() != 7252 {
		t.Errorf("got %d docs, expected %d\n", idx.Len(), 7252)
	}
}

func TestOpenIndex(t *testing.T) {
	params, err := param.Parse(cfgStr)
	if err != nil {
		t.Error(err)
	}

	idx := Open(params.SrchCfg)
	docs := idx.Len()
	if docs != 7252 {
		t.Errorf("got %d docs, expected %d\n", docs, 7252)
	}
}

func TestBleveSearch(t *testing.T) {
	params, err := param.Parse(cfgStr)
	if err != nil {
		t.Error(err)
	}

	idx := Open(params.SrchCfg)
	bits, err := idx.Search("fish")
	if err != nil {
		t.Error(err)
	}
	if h := bits.GetCardinality(); h != 8 {
		t.Errorf("got %d hits, expected %d\n", h, 8)
	}

	books := loadData(t)
	for id, doc := range books {
		if i, ok := doc["id"]; ok {
			id = cast.ToInt(i)
		}
		if bits.ContainsInt(id) {
			if title, ok := doc["title"].(string); ok {
				if !strings.Contains(strings.ToLower(title), "fish") {
					t.Errorf("result %s, doesn't contain query %s\n", title, "fish")
				}
			} else {
				t.Errorf("no field\n")
			}
		}
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
		books = append(books, map[string]any{
			"title": b["title"],
			//"id":    b["id"],
		})
	}

	return books
}
