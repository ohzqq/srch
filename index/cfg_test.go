package index

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

type cfgTest struct {
	params
	*param.Cfg
}

var cfgTests = []cfgTest{
	cfgTest{
		params: params{
			query: ``,
		},
		Cfg: &param.Cfg{
			SrchAttr: []string{"*"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		params: params{
			query: `?searchableAttributes=`,
		},
		Cfg: &param.Cfg{
			SrchAttr: []string{"*"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		params: params{
			query: `?searchableAttributes=title`,
		},
		Cfg: &param.Cfg{
			SrchAttr: []string{"title"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		params: params{
			query: `?searchableAttributes=title&url=file://../testdata/data-dir&sortableAttributes=tags`,
		},
		Cfg: &param.Cfg{
			SrchAttr: []string{"title"},
			Paramz: &param.Paramz{
				Index: "default",
				URI:   `file://../testdata/data-dir`,
			},
			SortAttr: []string{"tags"},
		},
	},
	cfgTest{
		params: params{
			query: `?attributesForFaceting=tags,authors,series,narrators`,
		},
		Cfg: &param.Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		params: params{
			query: `?attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`,
		},
		Cfg: &param.Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		params: params{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		},
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &param.Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		params: params{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&id=id`,
		},
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &param.Paramz{
				Index: "default",
				UID:   "id",
			},
		},
	},
	cfgTest{
		params: params{
			query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&url=file://../testdata/data-dir/audiobooks.json&index=audiobooks&id=id`,
		},
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			SortAttr:  []string{"tags"},
			Paramz: &param.Paramz{
				UID:   "id",
				Index: "audiobooks",
				URI:   "file://../testdata/data-dir/audiobooks.json",
			},
		},
	},
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
