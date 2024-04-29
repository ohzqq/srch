package data

import (
	"testing"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/srch/param"
)

const hareTestDB = `testdata/hare`

func TestNewDB(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		t.Error(err)
	}

	if !db.TableExists("index") {
		t.Error("table doesn't exist")
	}
}

func TestInsertRecords(t *testing.T) {
	t.SkipNow()
	//db := newDB()
	d, err := newData()
	if err != nil {
		t.Error(err)
	}

	params := param.New()
	params.SrchAttr = []string{"title", "comments"}
	params.Facets = []string{"tags"}

	var docs []*Doc
	for _, da := range d.data {
		docs = append(docs, NewDoc(da, params))
	}

	ds, err := disk.New(hareTestDB, ".json")
	if err != nil {
		t.Error(err)
	}
	db, err := hare.New(ds)
	if err != nil {
		t.Error(err)
	}

	err = db.CreateTable("index")
	if err != nil {
		t.Error(err)
	}

	for _, doc := range docs {
		id, err := db.Insert("index", doc)
		if err != nil {
			t.Error(err)
		}
		if id != doc.ID {
			t.Errorf("got id %v, expected %v\n", id, doc.ID)
		}
	}

	//err = db.DropTable("index")
	//if err != nil {
	//  t.Error(err)
	//}

}

func newDB() *DB {
	db, _ := NewDB()
	return db
}
