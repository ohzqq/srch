package mem

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/ohzqq/srch/param"
)

var testQuerySettings = []string{
	"searchableAttributes=title&dataDir=../testdata/nddata",
	"searchableAttributes=title&dataFile=../testdata/nddata/ndbooks.ndjson",
	"searchableAttributes=*&dataDir=../testdata/nddata",
	"",
	"searchableAttributes=",
}

func TestNewIndex(t *testing.T) {
	for i := 0; i < len(testQuerySettings); i++ {
		q := testQuerySettings[i]

		params, err := param.Parse(q)
		if err != nil {
			t.Error(err)
		}

		idx := New(params.SrchCfg)

		data := loadData(t)
		err = idx.Batch(data)
		if err != nil {
			t.Fatal(err)
		}

		total := idx.Len()
		switch i {
		case 0, 1, 2:
			if total != 7252 {
				t.Errorf("got %d, expected %d\n", total, 7252)
			}
		default:
			if total != 0 {
				t.Errorf("got %d, expected %d\n", total, 0)
			}
		}
	}
}

func TestSearchMem(t *testing.T) {
	for i := 0; i < len(testQuerySettings); i++ {
		q := testQuerySettings[i]

		params, err := param.Parse(q)
		if err != nil {
			t.Error(err)
		}

		idx := Open(params.SrchCfg)

		data := loadData(t)
		err = idx.Batch(data)
		if err != nil {
			t.Fatal(err)
		}

		res, err := idx.Search("fish")
		if err != nil {
			t.Error(err)
		}
		total := len(res)

		switch i {
		case 2, 3:
			if total != 56 {
				t.Errorf("got %d, expected %d\n", total, 56)
			}
		default:
			if total != 0 {
				t.Errorf("got %d, expected %d\n", total, 0)
			}
		}
	}
}

func loadData(t *testing.T) []map[string]any {
	d, err := os.Open("../testdata/nddata/ndbooks.ndjson")
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
