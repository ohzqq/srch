package index

import "testing"

func TestDefaultIndex(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Error(err)
	}
}

func TestSettings(t *testing.T) {
	idx, err := New()
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
	idx, err := New()
	if err != nil {
		t.Fatal(err)
	}

	tbl, err := idx.Cfg()
	if err != nil {
		t.Error(err)
	}

	if len(tbl) != 1 {
		t.Errorf("got %v, wanted %v\n", len(tbl), 1)
	}
}

func TestGetCfg(t *testing.T) {
	idx, err := New()
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := idx.GetCfg(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if cfg.Name != defaultTbl {
		t.Errorf("got %v, wanted %v\n", cfg.Name, defaultTbl)
	}
}
