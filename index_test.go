package srch

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/mitchellh/mapstructure"
)

var idx = &Index{}

var books []any

const numBooks = 7174

const testData = `testdata/data-dir/audiobooks.json`
const testCfgFile = `testdata/config.json`
const testCfgFileData = `testdata/config-with-data.json`

func init() {
	var err error
	idx, err = New(testCfgFile,
		DataFile(testData),
		//WithSearch(testS),
	)
	if err != nil {
		log.Fatal(err)
	}
	books = idx.Data

}

func TestIdxCfg(t *testing.T) {
	//cfg := &Index{}
	err := json.Unmarshal([]byte(testCfg), idx)
	if err != nil {
		t.Error(err)
	}
}

func TestIdxFilterSearch(t *testing.T) {
	t.SkipNow()
	//vals := testVals()
	//res := idx.Search(vals)

	fn := FuzzySearch(books, "title")
	res := fn("fish")
	i, err := New(testCfgFile, DataSlice(res))
	if err != nil {
		t.Error(err)
	}
	vals := make(url.Values)
	vals.Set("authors", "amy lane")
	r := i.Filter(vals)
	if len(r.Data) != 4 {
		t.Errorf("got %d, expected 4", len(r.Data))
	}
}

func TestIdxSearch(t *testing.T) {
	println("test idx search")
	vals := testVals()
	r := idx.Search(vals)
	//fmt.Println(len(r.Data))
	r.Print()
}

func TestNewIdxFromString(t *testing.T) {
	idx, err := parseCfg(testCfg)
	if err != nil {
		t.Error(err)
	}
	if len(idx.Facets) != 2 {
		t.Errorf("got %d facets, expected 2", len(idx.Facets))
	}

	d, err := os.ReadFile(testData)
	if err != nil {
		t.Error(err)
	}

	data, err := NewDataFromString(string(d))
	if err != nil {
		t.Error(err)
	}
	if len(data) != len(books) {
		t.Errorf("got %d, expected 7174\v", len(data))
	}
}

func TestNewIdxFromMap(t *testing.T) {
	t.SkipNow()
	d := make(map[string]any)
	err := mapstructure.Decode(idx, &d)
	if err != nil {
		t.Error(err)
	}
	i, err := NewIndexFromMap(d)
	if err != nil {
		t.Error(err)
	}
	if len(i.Data) != len(books) {
		t.Errorf("got %d, expected 7174\v", len(i.Data))
	}
	if len(i.Facets) != 2 {
		t.Errorf("got %d facets, expected 2", len(i.Facets))
	}
}

func loadData(t *testing.T) []map[string]any {
	d, err := os.ReadFile(testData)
	if err != nil {
		t.Error(err)
	}

	var books []map[string]any
	err = json.Unmarshal(d, &books)
	if err != nil {
		t.Error(err)
	}

	books = books

	return books
}

func TestData(t *testing.T) {
	books := loadData(t)
	if len(books) != 7174 {
		t.Errorf("got %d, expected 7174\v", len(books))
	}
}

const testCfg = `{
	"facets": [
		{
			"attribute": "tags",
			"operator": "and"
		},
		{
			"attribute": "authors"
		}
	]
}
`
