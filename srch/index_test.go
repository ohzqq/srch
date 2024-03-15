package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ohzqq/srch/param"
	"github.com/ohzqq/srch/txt"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7253

const testDataFile = `../testdata/data-dir/audiobooks.json`
const testDataDir = `../testdata/data-dir`
const testCfgFile = `../testdata/config-old.json`
const testYAMLCfgFile = `../testdata/config.yaml`
const testCfgFileData = `../testdata/config-with-data.json`
const libCfgStr = "searchableAttributes=title&facets=tags,authors,series,narrators&dataFile=testdata/data-dir/audiobooks.json"

func TestData(t *testing.T) {
	books = loadData(t)
	err := totalBooksErr(len(books), 71734)
	if err != nil {
		t.Error(err)
	}
}

var testQueryNewIndex = []string{
	"searchableAttributes=title&fullText",
	"searchableAttributes=title&dataDir=../testdata/data-dir",
	"facets=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&facets=tags,authors,series,narrators",
}

var testQuerySettings = []string{
	"",
	"searchableAttributes=",
	"searchableAttributes=title&fullText=../testdata/poot.bleve",
	"searchableAttributes=title&dataDir=../testdata/data-dir",
	"facets=tags,authors,series,narrators",
	"facets=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&facets=tags,authors,series,narrators",
	"searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&facets=tags,authors,series,narrators",
	`searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&facets=tags,authors,series,narrators&page=3&query=fish&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
}

var titleField = txt.NewField(param.DefaultField)

func TestNewIndex(t *testing.T) {
	for i := 0; i < len(testQuerySettings); i++ {
		q := testQuerySettings[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		var num int
		if !idx.Params.HasData() {
			data := loadData(t)
			num = len(data)
		} else {
			//idx.Index(data)
			num = idx.NbHits()
		}
		err = totalBooksErr(num, q)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestNewIndexWithParams(t *testing.T) {
	for i := 0; i < len(testQuerySettings); i++ {
		q := testQuerySettings[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		var num int
		if !idx.Params.HasData() {
			data := loadData(t)
			num = len(data)
		} else {
			//idx.Index(data)
			//num = 0
			num = idx.NbHits()
		}
		err = totalBooksErr(num, q)
		if err != nil {
			t.Error(err)
		}
	}
}

func intErr(got, want int, msg ...any) error {
	if got != want {
		err := fmt.Errorf("got %d, want %d\n", got, want)
		if len(msg) > 0 {
			err = fmt.Errorf("%w [msg] %v\n", err, msg)
		}
		return err
	}
	return nil
}

func totalBooksTest(total int, t *testing.T) {
	if total != numBooks {
		t.Errorf("got %d, expected %d\n", total, numBooks)
	}
}

func newTestIdx() *Index {
	//idx, err := New(libCfgStr)
	//if err != nil {
	//log.Fatal(err)
	//}
	return newTestIdxCfg("")
}

func newTestIdxCfg(p string) *Index {
	idx, err := New(libCfgStr + p)
	if err != nil {
		log.Fatal(err)
	}
	return idx
}

func totalBooksErr(total int, vals ...any) error {
	if total != numBooks {
		err := fmt.Errorf("got %d, expected %d\n", total, numBooks)
		return fmt.Errorf("%w\nmsg: %v", err, vals)
	}
	return nil
}

func loadData(t *testing.T) []map[string]any {
	d, err := os.ReadFile(testDataFile)
	if err != nil {
		t.Error(err)
	}

	var books []map[string]any
	err = json.Unmarshal(d, &books)
	if err != nil {
		t.Error(err)
	}

	return books
}
