package srch

import (
	"encoding/json"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

//func TestDocMapping(t *testing.T) {
//  doc := DocMapping(testFieldMappings())
//  d, err := json.Marshal(doc)
//  if err != nil {
//    t.Error(err)
//  }
//  println(string(d))
//}

func TestNewBleveIndex(t *testing.T) {
	idxMap := bleve.NewIndexMapping()
	idxMap.DefaultAnalyzer = StandardAnalyzer

	title := bleve.NewTextFieldMapping()
	title.Name = "title"

	id := bleve.NewTextFieldMapping()
	id.Name = "id"

	//books := bleve.NewDocumentStaticMapping()
	//books.AddFieldMapping(title)

	idxMap.DefaultMapping = bleve.NewDocumentStaticMapping()
	idxMap.DefaultMapping.DefaultAnalyzer = StandardAnalyzer
	//idxMap.DefaultMapping.AddFieldMapping(title)
	idxMap.DefaultMapping.AddFieldMappingsAt("title", title)
	idxMap.DefaultMapping.AddFieldMappingsAt("id", id)

	//idxMap.AddDocumentMapping("title", books)

	idx, err := bleve.New("testdata/poot", idxMap)
	if err != nil {
		t.Error(err)
	}

	d, err := json.Marshal(idx)
	if err != nil {
		t.Error(err)
	}
	println(string(d))
}

func testFieldMappings() []*mapping.FieldMapping {
	fields := []string{DefaultField}
	return FieldMappings(fields)
}
