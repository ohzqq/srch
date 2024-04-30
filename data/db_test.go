package data

import (
	"testing"

	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
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

	if len(ids) != 2 {
		t.Errorf("got %v results, expected %v\n", len(ids), 2)
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
	db := newDB()
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

	for _, doc := range docs {
		id, err := db.Insert(doc)
		if err != nil {
			t.Error(err)
		}
		if id != doc.ID {
			t.Errorf("got id %v, expected %v\n", id, doc.ID)
		}
	}

	err = db.DropTable("index")
	if err != nil {
		t.Error(err)
	}
}

func TestInsertRecordsDisk(t *testing.T) {
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

	db, err := NewDiskDB(hareTestDB)
	if err != nil {
		t.Error(err)
	}

	for _, doc := range docs {
		id, err := db.Insert(doc)
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
