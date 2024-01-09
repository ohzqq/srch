package srch

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mitchellh/mapstructure"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7174

const testData = `testdata/data-dir/audiobooks.json`
const testCfgFile = `testdata/config.json`
const testCfgFileData = `testdata/config-with-data.json`
const testQueryString = `tags=grumpy/sunshine&tags=enemies+to+lovers`

func init() {
	idx = New(DataFile(testData), WithCfg(testCfgFile))
}

func TestData(t *testing.T) {
	books = loadData(t)
	if len(books) != 7174 {
		t.Errorf("got %d, expected 7174\v", len(books))
	}
}

func TestNewIndex(t *testing.T) {
	if idx.Len() != len(books) {
		t.Errorf("got %d, expected %d\n", idx.Len(), len(books))
	}
}

func TestIndexProps(t *testing.T) {
	if len(idx.Facets()) != 4 {
		t.Errorf("got %d, expected 4\n", len(idx.Facets()))
	}
	if len(idx.TextFields()) != 1 {
		t.Errorf("got %d, expected 4\n", len(idx.TextFields()))
	}
}

func TestIdxCfgString(t *testing.T) {
	istr := New(CfgString(testCfg))
	facets := FilterFacets(istr.Fields)
	if len(istr.Facets()) != len(facets) {
		t.Errorf("got %d, expected %d\n", len(istr.Facets()), len(facets))
	}
	if len(istr.TextFields()) != 1 {
		t.Errorf("got %d, expected 4\n", len(istr.TextFields()))
	}
}

func TestNewIdxFromMap(t *testing.T) {
	t.SkipNow()
	d := make(map[string]any)
	err := mapstructure.Decode(idx, &d)
	if err != nil {
		t.Error(err)
	}
	err = CfgIndexFromMap(idx, d)
	if err != nil {
		t.Error(err)
	}
	data := idx.GetData()
	if len(data) != len(books) {
		t.Errorf("got %d, expected 7174\v", len(data))
	}
	if len(idx.Facets()) != 2 {
		t.Errorf("got %d facets, expected 2", len(idx.Facets()))
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
