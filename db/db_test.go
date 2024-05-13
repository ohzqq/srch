package db

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

const hareTestDB = `testdata/hare`

func TestNewDB(t *testing.T) {
	//cleanFiles(t)

	db, err := New(NewDisk(hareTestDB))
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"default"}

	if !slices.Equal(db.ListTables(), want) {
		t.Errorf("got %v tables, wanted %v\n", db.ListTables(), want)
	}

	if db.TableExists(defaultTbl) {
		err = db.DropTable(defaultTbl)
		if err != nil {
			t.Error(err)
		}

		err = db.CreateTable(defaultTbl)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCfgTable(t *testing.T) {
	m := testMapping()
	db, err := New(
		WithDisk(hareTestDB),
		WithDefaultCfg(defaultTbl, m, "id"),
	)
	if err != nil {
		t.Error(err)
	}

	cfg, err := db.GetTable(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if cfg.Name != defaultTbl {
		t.Errorf("wanted name %v, got %v\n", defaultTbl, cfg.Name)
	}
}

func TestOpenDB(t *testing.T) {
	db, err := openDiskDB()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		defaultTbl,
	}
	if !slices.Equal(db.ListTables(), want) {
		t.Errorf("got %v tables, wanted %v\n", db.TableNames(), want)
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

	if db.TableExists(defaultTbl) {
		err = db.DropTable(defaultTbl)
		if err != nil {
			t.Error(err)
		}

		err = db.CreateTable(defaultTbl)
		if err != nil {
			t.Error(err)
		}
	}

	err = batchInsert(db)
	if err != nil {
		t.Error(err)
	}
}

func TestMemHare(t *testing.T) {
	t.SkipNow()
	tp := filepath.Join(hareTestDB, "default.json")

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

	tbl, err := db.GetTable(defaultTbl)
	if err != nil {
		t.Error(err)
	}

	ids, err := tbl.IDs()
	if err != nil {
		t.Error(err)
	}

	if len(ids) != 7251 {
		t.Errorf("got %v, wanted %v\n", len(ids), 7251)
	}

}

func TestAllRecs(t *testing.T) {
	//t.SkipNow()
	//dsk, err := NewDiskStorage(hareTestDB)
	dsk, err := disk.New(hareTestDB, ".json")
	if err != nil {
		t.Error(err)
	}
	h, err := hare.New(dsk)
	if err != nil {
		t.Error(err)
	}

	tb, err := h.GetTable(defaultTbl)
	if err != nil {
		t.Error(err)
	}

	ids, err := tb.IDs()
	if err != nil {
		t.Error(err)
	}

	db, err := New(WithDisk(hareTestDB))
	if err != nil {
		t.Error(err)
	}
	tbl, err := db.GetTable(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	res, err := tbl.IDs()
	if err != nil {
		t.Error(err)
	}

	if len(res) != len(ids) {
		t.Errorf("got %v, want %v\n", len(res), len(ids))
	}
}

func TestFindRec(t *testing.T) {
	//t.SkipNow()

	db, err := New(WithDisk(hareTestDB))
	if err != nil {
		t.Fatal(err)
	}
	tbl, err := db.GetTable(defaultTbl)
	if err != nil {
		t.Fatal(err)
	}
	find := 1832
	_, err = tbl.Find(find)
	if err != nil {
		t.Error(err)
	}
}

func TestNewRamDB(t *testing.T) {
	//t.SkipNow()

	data, err := os.ReadFile("/home/mxb/code/srch/testdata/hare/default.json")
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(WithData(data))
	if err != nil {
		t.Error(err)
	}

	ids, err := db.IDs(defaultTbl)
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

	ids, err := db.IDs(defaultTbl)
	if err != nil {
		t.Error(err)
	}

	if len(ids) != 7251 {
		t.Errorf("got %v, want %v\n", len(ids), 7251)
	}
}

func TestNewNet(t *testing.T) {
	d, err := os.ReadFile(filepath.Join(hareTestDB, "default.json"))
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(WithURL("http://mxb.ca/search/default.json", d))
	if err != nil {
		t.Fatal(err)
	}

	total, err := db.IDs(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if len(total) != 7251 {
		t.Errorf("got %v, wanted %v\n", len(total), 7251)
	}
}

func TestCreateTable(t *testing.T) {
}

func openDiskDB() (*Client, error) {
	db, err := New(
		WithDisk(hareTestDB),
	)
	return db, err
}

func testMapping() doc.Mapping {
	m := doc.NewMapping()
	m.AddFulltext("title", "comments")
	m.AddKeywords("tags", "authors", "narrators", "series")
	return m
}

func testParams() *param.Params {
	params := param.New()
	params.UID = "id"
	//params.SrchAttr = []string{"title"}
	//params.SrchAttr = []string{"comments"}
	params.SrchAttr = []string{"title", "comments"}
	params.Facets = []string{"tags", "authors", "narrators", "series"}
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

func batchInsert(db *Client) error {
	data, err := os.ReadFile("/home/mxb/code/srch/testdata/hare/default.json")
	if err != nil {
		return err
	}

	tbl, err := db.GetTable(defaultTbl)
	if err != nil {
		return err
	}

	err = tbl.Batch(data)
	if err != nil {
		return err
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
	"default.json",
}
