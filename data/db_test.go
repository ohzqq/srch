package data

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

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
	db := newDB()
	d, err := newData()
	if err != nil {
		t.Error(err)
	}

	params := param.New()
	params.SrchAttr = []string{"title"}
	params.Facets = []string{"tags"}

	for _, da := range d.data {
		doc := NewDoc(da, params)
		id, err := db.Database.Insert("index", doc)
		if err != nil {
			t.Error(err)
		}
		if id != doc.ID {
			t.Errorf("got id %v, expected %v\n", id, doc.ID)
		}
	}
}

func newDB() *DB {
	db, _ := NewDB()
	return db
}
