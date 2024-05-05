package db

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

const hareTestDB = `testdata/hare`

func TestMemHare(t *testing.T) {
	tp := filepath.Join(hareTestDB, "index.json")

	data := data.New("file", tp)
	d, err := data.Docs()
	if err != nil {
		t.Error(err)
	}

	ds, err := ram.New(d)
	if err != nil {
		t.Error(err)
	}

	db, err := hare.New(ds)
	if err != nil {
		t.Error(err)
	}

	ids, err := db.IDs("index")
	if err != nil {
		t.Error(err)
	}

	if len(ids) != 7251 {
		t.Errorf("got %v, wanted %v\n", len(ids), 7251)
	}

}

func TestAllRecs(t *testing.T) {
	//t.SkipNow()
	dsk, err := NewDisk(hareTestDB)
	if err != nil {
		t.Error(err)
	}

	ids, err := dsk.IDs("index")
	if err != nil {
		t.Error(err)
	}
	slices.Sort(ids)
	i := 1
	for _, id := range ids {
		if id != i {
			println(i)
		}
		i++
	}

	db, err := New(WithDisk(hareTestDB))
	if err != nil {
		t.Error(err)
	}
	res, err := db.Find(-1)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 7251 {
		t.Errorf("got %v, want %v\n", len(res), 7251)
	}
}

func TestFindRec(t *testing.T) {
	//t.SkipNow()

	db, err := New(WithDisk(hareTestDB))
	if err != nil {
		t.Fatal(err)
	}
	find := 1832
	_, err = db.Find(find)
	if err != nil {
		t.Error(err)
	}
	//found := doc.SearchAllFields("range")
	//if !found {
	//t.Errorf("%#v\n", doc)
	//}
}

func TestNewRamDB(t *testing.T) {
	//t.SkipNow()

	data, err := os.ReadFile(`../testdata/ndbooks.ndjson`)
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(WithData(data))
	if err != nil {
		t.Error(err)
	}

	ids, err := db.IDs("index")
	if err != nil {
		t.Error(err)
	}

	if len(ids) != 7251 {
		t.Errorf("got %v, want %v\n", len(ids), 7251)
	}
}

func TestInsertRamDB(t *testing.T) {
	//t.SkipNow()

	data, err := os.ReadFile(`../testdata/ndbooks.ndjson`)
	if err != nil {
		t.Fatal(err)
	}

	db, err := New()
	if err != nil {
		t.Error(err)
	}

	err = db.Batch(testMapping(), data)
	if err != nil {
		t.Error(err)
	}

	ids, err := db.IDs("index")
	if err != nil {
		t.Error(err)
	}

	if len(ids) != 7251 {
		t.Errorf("got %v, want %v\n", len(ids), 7251)
	}

}

func TestNewNet(t *testing.T) {
	d, err := os.ReadFile(`../testdata/ndbooks.ndjson`)
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(WithURL("http://mxb.ca/search/index.json", d))
	if err != nil {
		t.Fatal(err)
	}

	total, err := db.IDs("index")
	if err != nil {
		t.Error(err)
	}
	if len(total) != 7251 {
		t.Errorf("got %v, wanted %v\n", len(total), 7251)
	}
}

func TestInsertRecordsDisk(t *testing.T) {
	//t.SkipNow()
	db, err := New(WithDisk(hareTestDB))
	if err != nil {
		t.Error(err)
	}

	if db.TableExists("index") {
		err = db.DropTable("index")
		if err != nil {
			t.Error(err)
		}
		err = db.CreateTable("index")
		if err != nil {
			t.Error(err)
		}
	}

	data, err := newData()
	if err != nil {
		t.Error(err)
	}

	err = db.BatchInsert(testMapping(), data)
	if err != nil {
		t.Error(err)
	}
}

func testMapping() *doc.Mapping {
	m := doc.NewMapping()
	m.AddFulltext("title", "comments", "tags")
	m.AddKeywords("tags")
	return m
}

func testParams() *param.Params {
	params := param.New()
	//params.UID = "id"
	//params.SrchAttr = []string{"title"}
	//params.SrchAttr = []string{"comments"}
	params.SrchAttr = []string{"title", "comments", "tags"}
	params.Facets = []string{"tags"}
	return params
}

func newData() ([]map[string]any, error) {
	d := data.New("file", `../testdata/ndbooks.ndjson`)

	recs, err := d.Decode()
	if err != nil {
		return nil, err
	}

	return recs, err
}
