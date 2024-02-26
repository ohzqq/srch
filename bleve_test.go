package srch

import (
	"encoding/json"
	"testing"

	"github.com/blevesearch/bleve/v2"
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
