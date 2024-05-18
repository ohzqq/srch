package index

import (
	"testing"
)

//func TestClientCfg(t *testing.T) {
//  for _, test := range cfgTests {
//    client, err := New(test.str())
//    if err != nil {
//      t.Fatal(err)
//    }
//  }
//}

func TestSettings(t *testing.T) {
	client, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Cfg()
	if err != nil {
		t.Error(err)
	}

	if !client.TableExists(settingsTbl) {
		t.Errorf("_settings table doesn't exist")
	}
}

func TestGetCfg(t *testing.T) {
	t.SkipNow()
	client, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := client.Cfg()
	if err != nil {
		t.Error(err)
	}

	idx, err := cfg.Find(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if idx.Index != defaultTbl {
		t.Errorf("got %v, wanted %v\n", idx.Index, defaultTbl)
	}
}
