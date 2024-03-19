package srch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7253

const testDataFile = `file/testdata/data-dir/audiobooks.json`
const testDataNdFile = `file/testdata/nddata/ndbooks.ndjson`
const testDataDir = `dir/testdata/data-dir`
const testBlvPath = `blv/testdata/poot.bleve`

var libCfgStr = mkURL(testDataNdFile, "searchableAttributes=title&facets=tags,authors,series,narrators")

func TestData(t *testing.T) {
	books = loadData(t)
	err := totalBooksErr(len(books), 71734)
	if err != nil {
		t.Error(err)
	}
}

var testQuerySettings = []string{
	"",
	"searchableAttributes=",
	mkURL(testBlvPath, "searchableAttributes=title"),
	mkURL(testDataDir, "searchableAttributes=title"),
	mkURL("", "facets=tags,authors,series,narrators"),
	mkURL(testDataFile, "facets=tags,authors,series,narrators"),
	mkURL(testDataDir, "searchableAttributes=title"),
	mkURL(testDataDir, "searchableAttributes=title"),
	mkURL("", "searchableAttributes=title&facets=tags,authors,series,narrators"),
	mkURL(testDataFile, "searchableAttributes=title&facets=tags,authors,series,narrators"),
	mkURL("", "searchableAttributes=title&facets=tags,authors,series,narrators"),
	mkURL(testDataNdFile, `searchableAttributes=title&facets=tags,authors,series,narrators&page=3&query=fish&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`),
}

func TestNewIndex(t *testing.T) {
	for i := 0; i < len(testQuerySettings); i++ {
		q := testQuerySettings[i]
		println(q)
		idx, err := New(q)
		if err != nil {
			if errors.Is(err, NoDataErr) {
				t.Log(err)
			} else {
				t.Error(err)
			}
		}
		var num int
		if !idx.Params.HasData() {
			data := loadData(t)
			num = len(data)
		} else {
			//idx.Index(data)
			num = idx.Len()
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
			if errors.Is(err, NoDataErr) {
				t.Log(err)
			} else {
				t.Error(err)
			}
		}
		var num int
		if !idx.Params.HasData() {
			data := loadData(t)
			num = len(data)
		} else {
			//idx.Index(data)
			//num = 0
			num = idx.Len()
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

func mkURL(path, rq string) string {
	u := &url.URL{
		Path:     path,
		RawQuery: rq,
	}
	return u.String()
}

func totalBooksErr(total int, vals ...any) error {
	if total != numBooks {
		err := fmt.Errorf("got %d, expected %d\n", total, numBooks)
		return fmt.Errorf("%w\nmsg: %v", err, vals)
	}
	return nil
}

func loadData(t *testing.T) []map[string]any {
	d, err := os.ReadFile(strings.TrimPrefix(testDataFile, "file/"))
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
