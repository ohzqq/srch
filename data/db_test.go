package data

import (
	"testing"

	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const hareTestDB = `testdata/hare`

func TestAllRecs(t *testing.T) {
	params := testParams()
	db, err := NewDB(params, WithHare(hareTestDB))
	if err != nil {
		t.Error(err)
	}
	res, err := db.AllRecords()
	if err != nil {
		t.Error(err)
	}
	if len(res) != 7252 {
		t.Errorf("got %v, want %v\n", len(res), 7252)
	}
}

func TestSearchDB(t *testing.T) {
	params := testParams()
	db, err := NewDB(params, WithHare(hareTestDB))
	if err != nil {
		t.Error(err)
	}

	ids, err := db.Search("falling fish")
	if err != nil {
		t.Error(err)
	}

	ids = lo.Uniq(ids)

	want := 10
	if len(ids) > want {
		t.Errorf("got %v results, expected %v\n", len(ids), want)
	}
}

func TestFindRec(t *testing.T) {
	params := testParams()
	db, err := NewDB(params, WithHare(hareTestDB))
	if err != nil {
		t.Error(err)
	}
	find := 1832
	doc, err := db.Find(find)
	if err != nil {
		t.Error(err)
	}
	found := doc.SearchAllFields("range")
	if !found {
		t.Errorf("%#v\n", doc)
	}
}

func TestInsertRecordsRam(t *testing.T) {
	//t.SkipNow()
	params := testParams()
	db, err := NewDB(params)
	if err != nil {
		t.Error(err)
	}

	d, err := newData()
	if err != nil {
		t.Error(err)
	}

	for id, dd := range d.data {
		doc, err := db.Insert(dd)
		if err != nil {
			t.Error(err)
		}
		if i, ok := dd[db.UID]; ok {
			id = cast.ToInt(i)
		}
		if doc.GetID() != id {
			t.Errorf("got id %v, expected %v\n", doc.GetID(), id)
		}
	}
}

func TestInsertRecordsDisk(t *testing.T) {
	t.SkipNow()
	//db := newDB()
	params := testParams()
	db, err := NewDB(params, WithHare(hareTestDB))
	if err != nil {
		t.Error(err)
	}

	d, err := newData()
	if err != nil {
		t.Error(err)
	}

	for id, dd := range d.data {
		doc, err := db.Insert(dd)
		if err != nil {
			t.Error(err)
		}
		if i, ok := dd[db.UID]; ok {
			id = cast.ToInt(i)
		}
		if doc.GetID() != id {
			t.Errorf("got id %v, expected %v\n", doc.GetID(), id)
		}
	}

	//err = db.DropTable("index")
	//if err != nil {
	//  t.Error(err)
	//}

}

func testParams() string {
	params := param.New()
	//params.UID = "id"
	params.SrchAttr = []string{"title", "comments"}
	params.Facets = []string{"tags"}
	return params.String()
}

func newDB() *DB {
	db, _ := NewDB(testParams())
	return db
}
