package index

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

var cfgTests = []test{
	test{
		query: ``,
		Cfg: &param.Cfg{
			SrchAttr: []string{"*"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title`,
		Cfg: &param.Cfg{
			SrchAttr: []string{"title"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title&url=file://home/mxb/code/srch/testdata/hare/&sortableAttributes=tags`,
		Cfg: &param.Cfg{
			SrchAttr: []string{"title"},
			Paramz: &param.Paramz{
				Index: "default",
				URI:   `file://home/mxb/code/srch/testdata/hare/`,
			},
			SortAttr: []string{"tags"},
		},
	},
	test{
		query: `?attributesForFaceting=tags,authors,series,narrators`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&uid=id`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &param.Paramz{
				Index: "default",
				UID:   "id",
			},
		},
	},
	test{
		query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&url=file://home/mxb/code/srch/testdata/hare/&index=audiobooks&uid=id`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			SortAttr:  []string{"tags"},
			Paramz: &param.Paramz{
				UID:   "id",
				Index: "audiobooks",
				URI:   "file://home/mxb/code/srch/testdata/hare/",
			},
		},
	},
}

func TestCfgClientParams(t *testing.T) {
	for _, test := range cfgTests {
		client, err := New(test.str())
		if err != nil {
			t.Fatal(err)
		}
		if !client.TableExists(settingsTbl) {
			t.Error(test.msg("_settings table doesn't exist"))
		}
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

	if !client.TableExists(settingsTbl) {
		t.Errorf("_settings table doesn't exist")
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
