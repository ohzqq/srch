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
	testDataFile = `file/testdata/nddata/ndbooks.ndjson`
	testDataDir  = `dir/testdata/data-dir`
	testBlvPath  = `blv/testdata/poot.bleve`
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

var libCfgStr = mkURL(testDataFile, "searchableAttributes=title&facets=tags,authors,series,narrators")

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
	blvRoute(srchAttrParam),
	dirRoute(srchAttrParam),
	mkURL("", facetParamSlice),
	fileRoute(facetParamStr),
	dirRoute(testDataDir, srchAttrParam),
	dirRoute(testDataDir, srchAttrParam),
	mkURL("", srchAttrParam, facetParamSlice),
	fileRoute(srchAttrParam, facetParamStr),
	mkURL("", srchAttrParam, facetParamStr),
	fileRoute(testDataFile, srchAttrParam, facetParamSlice, `&page=3`, queryParam, sortParam, filterParam),
}

func TestNewIndex(t *testing.T) {
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
				t.Log(err)
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
		log.Fatal(err)
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
	return mkURL(testBlvPath, params...)
}

func dirRoute(params ...string) string {
	return mkURL(testDataDir, params...)
}

func fileRoute(params ...string) string {
	return mkURL(testDataFile, params...)
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
