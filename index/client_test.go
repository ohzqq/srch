package index

import "testing"

func TestDefaultIndex(t *testing.T) {
	_, err := New("")
	if err != nil {
		t.Error(err)
	}
}

func TestSettings(t *testing.T) {
	idx, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	_, err = idx.Cfg()
	if err != nil {
		t.Error(err)
	}

	if !idx.TableExists(settingsTbl) {
		t.Errorf("settings doesn't exist")
	}
}

func TestDefaultSettings(t *testing.T) {
	idx, err := New("")
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
