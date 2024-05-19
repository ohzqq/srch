package index

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

func TestClientInit(t *testing.T) {
	for _, test := range cfgTests {
		cfg := param.NewClient()
		err := param.Decode(test.str(), cfg)
		if err != nil {
			t.Error(err)
		}
		if cfg.DB != "" {
			if cfg.Path != hareTestPath {
				t.Errorf("got %v, wanted %v\n", cfg.Path, hareTestPath)
			}
			if cfg.Scheme != "file" {
				t.Errorf("got %v, wanted %v\n", cfg.Scheme, "file")
			}
		}
	}
}

func TestClientInitStr(t *testing.T) {
	for _, test := range cfgTests {
		client, err := New(test.str())
		if err != nil {
			t.Fatal(err)
		}

		//if cfg.DB != "" {
		//  if cfg.Path != hareTestPath {
		//    t.Errorf("got %v, wanted %v\n", cfg.Path, hareTestPath)
		//  }
		//  if cfg.Scheme != "file" {
		//    t.Errorf("got %v, wanted %v\n", cfg.Scheme, "file")
		//  }
		//}

		if !client.TableExists(settingsTbl) {
			t.Error(test.msg("_settings table doesn't exist"))
		}
		_, err = client.GetIdxCfg(client.Client.Index)
		if err != nil {
			t.Error(test.err(client.Client.Index, test.Cfg.Index, err))
		}
	}
}

func TestClientInitURL(t *testing.T) {
	for _, test := range cfgTests {
		client, err := New(test.url())
		if err != nil {
			t.Fatal(err)
		}
		if !client.TableExists(settingsTbl) {
			t.Error(test.msg("_settings table doesn't exist"))
		}
		_, err = client.GetIdxCfg(client.Client.Index)
		if err != nil {
			t.Error(test.err(client.Client.Index, test.Cfg.Index, err))
		}
	}
}

func TestClientInitValues(t *testing.T) {
	for _, test := range cfgTests {
		client, err := New(test.vals())
		if err != nil {
			t.Fatal(err)
		}
		if !client.TableExists(settingsTbl) {
			t.Error(test.msg("_settings table doesn't exist"))
		}
		_, err = client.GetIdxCfg(client.Client.Index)
		if err != nil {
			t.Error(test.err(client.Client.Index, test.Cfg.Index, err))
		}
	}
}

func TestHareDisk(t *testing.T) {
	_, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultClient(t *testing.T) {
	_, err := New("")
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultSettings(t *testing.T) {
	idx, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := idx.Cfg()
	if err != nil {
		t.Error(err)
	}

	ids, err := cfg.tbl.IDs()
	if err != nil {
		t.Error(err)
	}

	if len(ids) != 1 {
		t.Errorf("got %v, wanted %v\n", len(ids), 1)
	}
}

func TestGetIdx(t *testing.T) {
	t.SkipNow()
	c, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	idx, err := c.GetIdx(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if idx.IndexName() != defaultTbl {
		t.Errorf("got %v, wanted %v\n", idx.IndexName(), defaultTbl)
	}
}
