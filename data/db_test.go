package data

import (
	"slices"
	"testing"

	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

const hareTestDB = `testdata/hare`

func TestAllRecs(t *testing.T) {
	//t.SkipNow()
	params := testParams()
	dsk, err := Open(hareTestDB)
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
	db, err := NewDB(dsk, m)
	if err != nil {
		t.Error(err)
	}
	res, err := db.Find(-1)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 7252 {
		t.Errorf("got %v, want %v\n", len(res), 7252)
	}
}

func TestSearchDB(t *testing.T) {
	params := testParams()
	dsk, err := Open(hareTestDB)
	if err != nil {
		t.Error(err)
	}

	m := doc.NewMappingFromParams(params)
	db, err := NewDB(dsk, m)
	if err != nil {
		t.Error(err)
	}

	//fmt.Printf("%#v\n", m)

	//ids, err := db.Search("falling love")
	ids, err := db.Search("dragon omega")
	if err != nil {
		t.Error(err)
	}

	//fmt.Printf("%#v\n", ids)
	want := 140
	if len(ids) > want {
		println(len(ids) > want)
		t.Errorf("got %v results, expected %v\n", len(ids), want)
	}
}

func TestFindRec(t *testing.T) {
	//t.SkipNow()
	params := testParams()
	dsk, err := Open(hareTestDB)
	if err != nil {
		t.Error(err)
	}

	m := doc.NewMappingFromParams(params)
	db, err := NewDB(dsk, m)
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

func TestInsertRecordsRam(t *testing.T) {
	//t.SkipNow()
	mem := NewMem()
	params := testParams()
	m := doc.NewMappingFromParams(params)
	db, err := NewDB(mem, m)
	if err != nil {
		t.Error(err)
	}

	data, err := newData()
	if err != nil {
		t.Error(err)
	}

	err = db.Batch(data.data)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertRecordsDisk(t *testing.T) {
	//t.SkipNow()
	params := testParams()
	dsk, err := NewDisk(hareTestDB)
	if err != nil {
		t.Fatal(err)
	}
	m := doc.NewMappingFromParams(params)

	db, err := NewDB(dsk, m)
	if err != nil {
		t.Error(err)
	}

	data, err := newData()
	if err != nil {
		t.Error(err)
	}

	err = db.Batch(data.data)
	if err != nil {
		t.Error(err)
	}

	//err = db.DropTable("index")
	//if err != nil {
	//  t.Error(err)
	//}

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
