package srch

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"

	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/param"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7252

const (
	testDataFile = `testdata/nddata/ndbooks.ndjson`
	testDataDir  = `testdata/data-dir`
	testBlvPath  = `testdata/poot.bleve`
)

const (
	facetParamStr   = `facets=tags,authors,series,narrators`
	facetParamSlice = `facets=tags&facets=authors&facets=series&facets=narrators`
	srchAttrParam   = "searchableAttributes=title"
	queryParam      = `query=fish`
	sortParam       = `sortBy=title&order=desc`
	filterParam     = `facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`
	uidParam        = `uid=id`
)

var libCfgStr = fileRoute("searchableAttributes=title&facets=tags,authors,series,narrators")

func TestData(t *testing.T) {
	books = loadData(t)
	err := totalBooksErr(len(books), 71734)
	if err != nil {
		t.Error(err)
	}
}

var testQuerySettings = []string{
	blvRoute(srchAttrParam),
	dirRoute(srchAttrParam),
	fileRoute(facetParamStr),
	dirRoute(srchAttrParam),
	dirRoute(srchAttrParam),
	fileRoute(srchAttrParam, facetParamStr),
	fileRoute(srchAttrParam, facetParamSlice, `&page=3`, queryParam, sortParam, filterParam),
}

func TestNewIndex(t *testing.T) {
	for i := 0; i < len(testQuerySettings); i++ {
		q := testQuerySettings[i]
		idx, err := New(q)
		if err != nil {
			t.Errorf("params %s\n ", q)
			if errors.Is(err, NoDataErr) {
				t.Errorf("test new index %v\n", err)
			} else {
				t.Error(err)
			}
		}
		var num int
		if idx.Params.Path == "" {
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
				t.Errorf("test new index %v\n", err)
			} else {
				t.Error(err)
			}
		}
		var num int
		if idx.Params.Path == "" {
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
		log.Fatalf("test new index %v\n", err)
	}
	return idx
}

func mkURL(path string, rq ...string) string {
	u := &url.URL{
		Path:     path,
		RawQuery: strings.Join(rq, "&"),
	}
	return u.String()
}

func blvRoute(params ...string) string {
	params = append(params, "path="+testBlvPath)
	return mkURL(param.Blv, params...)
}

func dirRoute(params ...string) string {
	params = append(params, "path="+testDataDir)
	return mkURL(param.Dir, params...)
}

func fileRoute(params ...string) string {
	params = append(params, "path="+testDataFile)
	return mkURL(param.File, params...)
}

func totalBooksErr(total int, vals ...any) error {
	if total != numBooks && total != 7253 {
		err := fmt.Errorf("got %d, expected %d\n", total, numBooks)
		return fmt.Errorf("%w\nmsg: %v", err, vals)
	}
	return nil
}

func loadData(t *testing.T) []map[string]any {
	d := data.New(param.File, `testdata/nddata/ndbooks.ndjson`)

	books, err := d.Decode()
	if err != nil {
		t.Error(err)
	}

	return books
}
