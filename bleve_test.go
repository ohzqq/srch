package srch

import (
	"encoding/json"
	"testing"

	"github.com/blevesearch/bleve/v2/mapping"
)

func TestFieldMappings(t *testing.T) {
	maps := testFieldMappings()
	d, err := json.Marshal(maps)
	if err != nil {
		t.Error(err)
	}
	println(string(d))
}

func TestDocMapping(t *testing.T) {
	doc := DocMapping(testFieldMappings())
	d, err := json.Marshal(doc)
	if err != nil {
		t.Error(err)
	}
	println(string(d))
}

func testFieldMappings() []*mapping.FieldMapping {
	fields := []string{DefaultField}
	return FieldMappings(fields)
}
