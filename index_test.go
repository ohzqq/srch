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
const testCfgFile = `testdata/config.json`
const testYAMLCfgFile = `testdata/config.yaml`
const testCfgFileData = `testdata/config-with-data.json`

func init() {
	idx = New(WithCfg(testCfgFile))
	idx.Index(FileSrc(testData))
	books = idx.Data
}

func TestData(t *testing.T) {
	books = loadData(t)
	if len(books) != 7174 {
		t.Errorf("got %d, expected 7174\v", len(books))
	}
}

//func TestNewIndex(t *testing.T) {
//  if idx.Len() != len(books) {
//    t.Errorf("got %d, expected %d\n", idx.Len(), len(books))
//  }
//}

func TestIndexProps(t *testing.T) {
	if len(idx.FacetFields()) != 4 {
		t.Errorf("got %d, expected 4\n", len(idx.FacetFields()))
	}
	if len(idx.TextFields()) != 1 {
		t.Errorf("got %d, expected 4\n", len(idx.TextFields()))
	}
}

func TestRecursiveSearch(t *testing.T) {
	t.SkipNow()
	idx.search = FullTextSrchFunc(idx.Data, idx.TextFields())
	res := idx.Search("fish")
	fmt.Printf("after search %d\n", len(res.Data))
	fmt.Printf("after search %+v\n", res.Facets()[0].Items[0])
}

func TestIdxCfgString(t *testing.T) {
	t.SkipNow()
	istr := New(CfgString(testCfg))
	facets := FilterFacets(istr.Fields)
	if len(istr.FacetFields()) != len(facets) {
		t.Errorf("got %d, expected %d\n", len(istr.FacetFields()), len(facets))
	}
	if len(istr.TextFields()) != 1 {
		t.Errorf("got %d, expected 4\n", len(istr.TextFields()))
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
