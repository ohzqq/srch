package index

import (
	"testing"
)

const hareTestPath = `/home/mxb/code/srch/testdata/hare`
const hareTestURL = `file://home/mxb/code/srch/testdata/hare`
const hareTestQuery = `?url=file://home/mxb/code/srch/testdata/hare/`

func TestHareDisk(t *testing.T) {
	_, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%#v\n", idx.Database)
}

func TestHareDiskTbls(t *testing.T) {
	client, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}
	names := client.TableNames()
	for _, n := range names {
		println(n)
	}
}

func TestDefaultIndex(t *testing.T) {
	_, err := New("")
	if err != nil {
		t.Error(err)
	}
}

func TestSettings(t *testing.T) {
	client, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Cfg()
	if err != nil {
		t.Error(err)
	}

	//if !client.Database.TableExists(settingsTbl) {
	//err = client.Database.CreateTable(settingsTbl)
	//if err != nil {
	//  t.Fatal(err)
	//}
	//err = client.SetCfg(DefaultCfg())
	//if err != nil {
	//  return err
	//}
	//}

	//err = client.initDB()
	//if err != nil {
	//  t.Fatal(err)
	//}

	//if !client.TableExists(settingsTbl) {
	//t.Errorf("settings doesn't exist")
	//}
}

func TestDefaultSettings(t *testing.T) {
	idx, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}

	tbl, err := idx.Cfg()
	if err != nil {
		t.Error(err)
	}

	ids, err := tbl.IDs()
	if err != nil {
		t.Error(err)
	}

	if len(ids) != 1 {
		t.Errorf("got %v, wanted %v\n", len(ids), 1)
	}
}

func TestGetCfg(t *testing.T) {
	idx, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := idx.GetCfg(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if cfg.Index != defaultTbl {
		t.Errorf("got %v, wanted %v\n", cfg.Index, defaultTbl)
	}
}

func TestGetIdx(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	idx, err := c.GetIdx(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if idx.Index != defaultTbl {
		t.Errorf("got %v, wanted %v\n", idx.Index, defaultTbl)
	}
}
