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

func init() {
	idx = New(DataFile(testData), WithCfg(testCfgFile))
	books = idx.GetData()
}

func TestNewIndex(t *testing.T) {
	if idx.Len() != 7174 {
		t.Errorf("got %d, expected 7174\n", idx.Len())
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

func TestIdxCfg(t *testing.T) {
	//cfg := &Index{}
	err := json.Unmarshal([]byte(testCfg), idx)
	if err != nil {
		t.Error(err)
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
