package db

import (
	"fmt"
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
	params := testParams()
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

	m := doc.NewMappingFromParams(params)
	db, err := New(m, WithDisk(hareTestDB))
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

func TestSearchDB(t *testing.T) {
	tests := []map[string][]string{
		map[string][]string{
			"dragon": []string{"title"},
		},
		map[string][]string{
			"omega": []string{"title"},
		},
		map[string][]string{
			"dragon omega": []string{"title"},
		},
		map[string][]string{
			"dragon": []string{"comments"},
		},
		map[string][]string{
			"omega": []string{"comments"},
		},
		map[string][]string{
			"dragon omega": []string{"comments"},
		},
		map[string][]string{
			"dragon": []string{"title", "comments"},
		},
		map[string][]string{
			"omega": []string{"title", "comments"},
		},
		map[string][]string{
			"dragon omega": []string{"title", "comments"},
		},
	}

	want := []int{
		74,
		97,
		1,
		185,
		328,
		23,
		200,
		345,
		23,
	}

	ds, err := NewDiskStore(hareTestDB)
	if err != nil {
		t.Fatal(err)
	}

	db, err := Open(ds)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		for kw, attrs := range test {
			params := param.New()
			params.SrchAttr = attrs

			//m := doc.NewMappingFromParams(params)
			//db, err := New(m)
			//if err != nil {
			//t.Error(err)
			//}
			//err = db.Init(dsk)
			//if err != nil {
			//t.Fatal(err)
			//}

			ids, err := db.Search(kw)
			if err != nil {
				t.Error(err)
			}
			if res := len(ids); res != want[i] {
				fmt.Printf("kw %s, attrs %v: %#v\n", kw, attrs, res)
				t.Errorf("got %v results, wanted %v\n", res, want[i])
			}
		}
	}
}

func TestFindRec(t *testing.T) {
	//t.SkipNow()
	params := testParams()
	m := doc.NewMappingFromParams(params)

	db, err := New(m, WithDisk(hareTestDB))
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

	params := testParams()
	m := doc.NewMappingFromParams(params)

	data, err := os.ReadFile(`../testdata/ndbooks.ndjson`)
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(m, WithData(data))
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

	//err = db.Batch(data)
	//if err != nil {
	//t.Error(err)
	//}
}

func TestInsertRamDB(t *testing.T) {
	//t.SkipNow()

	//params := testParams()
	//m := doc.NewMappingFromParams(params)
	m := doc.DefaultMapping()

	data, err := os.ReadFile(`../testdata/ndbooks.ndjson`)
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(m)
	if err != nil {
		t.Error(err)
	}

	err = db.Batch(data)
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

	params := testParams()
	m := doc.NewMappingFromParams(params)

	db, err := New(m, WithURL("http://mxb.ca/search/index.json", d))
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
	params := testParams()
	m := doc.NewMappingFromParams(params)

	db, err := New(m, WithDisk(hareTestDB))
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

	err = db.BatchInsert(data)
	if err != nil {
		t.Error(err)
	}
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
