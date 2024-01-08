package srch

import (
	"encoding/json"
	"fmt"
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
	//var err error
	//idx, err = New(testCfgFile,
	//DataFile(testData),
	////WithSearch(testS),
	//)
	//if err != nil {
	//  log.Fatal(err)
	//}

	idx = New(FileSrc(testData), WithCfg(testCfgFile))

	books = idx.GetData()

}

func TestNewIndexFunc(t *testing.T) {
	i := New(FileSrc(testData), CfgFile(testCfgFile))
	if i.Len() != 7174 {
		t.Errorf("got %d, expected 7174\n", i.Len())
	}
	if len(i.Facets()) != 4 {
		t.Errorf("got %d, expected 4\n", len(i.Facets()))
	}
	field, err := i.GetField("tags")
	if err != nil {
		t.Error(err)
	}
	res := field.Filter("abo", "dnr")
	fmt.Printf("%v\n", len(res.ToArray()))
	//for _, f := range idx.Fields {
	//  fmt.Printf("%#v\n", f)
	//}
}

func TestIndexProps(t *testing.T) {
	for _, f := range idx.Facets() {
		fmt.Printf("attr %s\n", f.Attribute)
	}
	for _, f := range idx.TextFields() {
		fmt.Printf("attr %s\n", f.Attribute)
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
