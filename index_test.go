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
const testDataDir = `testdata/data-dir`
const testCfgFile = `testdata/config-old.json`
const testYAMLCfgFile = `testdata/config.yaml`
const testCfgFileData = `testdata/config-with-data.json`

func TestData(t *testing.T) {
	books = loadData(t)
	err := totalBooksErr(len(books), 71734)
	if err != nil {
		t.Error(err)
	}
}

var testQueryNewIndex = []string{
	"searchableAttributes=title",
	"searchableAttributes=title&dataDir=testdata/data-dir",
	"attributesForFaceting=tags,authors,series,narrators&dataFile=testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&dataFile=testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators",
}

var titleField = NewField(DefaultField, Fuzzy)

func TestNewIndex(t *testing.T) {
	for i := 0; i < len(testQueryNewIndex); i++ {
		q := testQueryNewIndex[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		switch i {
		case 0:
			data := loadData(t)
			idx.Index(data)
		case 1:
		case 2:
		case 3:
		}
		err = totalBooksErr(idx.Len(), q)
		if err != nil {
			t.Error(err)
		}
	}
}

func totalBooksTest(total int, t *testing.T) {
	if total != 7174 {
		t.Errorf("got %d, expected %d\n", total, 7174)
	}
}

func totalBooksErr(total int, vals ...any) error {
	if total != 7174 {
		err := fmt.Errorf("got %d, expected %d\n", total, 7174)
		return fmt.Errorf("%w\nmsg: %v", err, vals)
	}
	return nil
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
