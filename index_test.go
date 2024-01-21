package srch

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7174

const testData = `testdata/data-dir/audiobooks.json`
const testCfgFile = `testdata/config-old.json`
const testYAMLCfgFile = `testdata/config.yaml`
const testCfgFileData = `testdata/config-with-data.json`

func init() {
	query := fmt.Sprintf("%s&%s&%s", testValuesCfg, testQueryString, testSearchString)
	idx = NewIndex(query)
	books = idx.Data
}

func TestData(t *testing.T) {
	books = loadData(t)
	totalBooksTest(len(books), t)
}

func totalBooksTest(total int, t *testing.T) {
	if total != 7174 {
		t.Errorf("got %d, expected %d\n", total, 7174)
	}
}

func TestNewIndex(t *testing.T) {
	data := loadData(t)
	for _, test := range settingsTestVals {
		idx := New(test.query).Index(data)
		totalBooksTest(idx.Len(), t)
	}
}

func TestSortIndex(t *testing.T) {
	q := getNewQuery()
	i := NewIndex(q.Encode())
	i.Sort()
	//for _, d := range i.Data {
	//  fmt.Printf("%s\n", d["title"])
	//}
}

func TestIndexProps(t *testing.T) {
	if len(idx.Facets()) != 4 {
		t.Errorf("got %d, expected 4\n", len(idx.Facets()))
	}
	if len(idx.TextFields()) != 1 {
		t.Errorf("got %d, expected %d\n", len(idx.TextFields()), 1)
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

const testCfg = `{
	"fields": [
		{
			"attribute": "title",
			"fieldType": "text",
			"operator": "and"
		},
		{ 
			"fieldType": "facet",
			"attribute": "series"
		},
		{
			"fieldType": "facet",
			"attribute": "tags",
			"operator": "and"
		}
	]
}
`
