package db

import (
	"bytes"
	"encoding/json"
	"io"
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

func TestNewDB(t *testing.T) {
	cleanFiles(t)

	db, err := New(NewDisk(hareTestDB))
	if err != nil {
		t.Fatal(err)
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
}

func TestInsertRecordsDisk(t *testing.T) {
	//t.SkipNow()
	//m := testMapping()
	db, err := New(
		WithDisk(hareTestDB),
	)
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

	err = batchInsert(db)
	if err != nil {
		t.Error(err)
	}
}

func TestOpenDB(t *testing.T) {
	db, err := openDiskDB()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"index",
	}
	if !slices.Equal(db.ListTables(), want) {
		t.Errorf("got %v tables, wanted %v\n", db.TableNames(), want)
	}
}

func TestMemHare(t *testing.T) {
	tp := filepath.Join(hareTestDB, "index.json")

	data := data.New("file", tp)
	d, err := data.Docs()
	if err != nil {
		t.Error(err)
	}

	ds, err := ram.NewWithTables(d)
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
	dsk, err := NewDiskStorage(hareTestDB)
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
	res, err := db.Find("index", -1)
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
	_, err = db.Find("index", find)
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

	data, err := os.ReadFile(filepath.Join(hareTestDB, "index.json"))
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

	db, err := New()
	if err != nil {
		t.Error(err)
	}

	err = batchInsert(db)
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
	d, err := os.ReadFile(filepath.Join(hareTestDB, "index.json"))
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

func TestCfgTable(t *testing.T) {
	m := testMapping()
	db, err := New(
		WithDisk(hareTestDB),
		WithDefaultCfg("index", m, "id"),
	)
	if err != nil {
		t.Error(err)
	}

	cfg, err := db.GetCfg("index")
	if err != nil {
		t.Error(err)
	}
	if cfg.Name != "index" {
		t.Errorf("wanted name %v, got %v\n", "index", cfg.Name)
	}
}

func TestCreateTable(t *testing.T) {
}

func openDiskDB() (*DB, error) {
	db, err := New(
		WithDisk(hareTestDB),
	)
	return db, err
}

func testMapping() doc.Mapping {
	m := doc.NewMapping()
	m.AddFulltext("title", "comments", "tags")
	m.AddKeywords("tags")
	return m
}

func testParams() *param.Params {
	params := param.New()
	params.UID = "id"
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

func batchInsert(db *DB) error {
	data, err := os.ReadFile("/home/mxb/code/srch/param/testdata/hare/index.json")
	if err != nil {
		return err
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	for {
		doc := &doc.Doc{}
		if err := dec.Decode(doc); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		_, err := db.Insert("index", doc)
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanFiles(t *testing.T) {
	for _, f := range testFiles {
		n := filepath.Join(hareTestDB, f)
		err := os.Remove(n)
		if err != nil {
			t.Error(err)
		}
	}
}

var testFiles = []string{
	"_settings.json",
	"index.json",
}
