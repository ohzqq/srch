package param

import (
	"fmt"
	"slices"
	"testing"
)

type cfgTest struct {
	pt
	*Cfg
}

var cfgTests = []cfgTest{
	cfgTest{
		pt: pt{
			query: ``,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"*"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"*"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"title"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&db=file://../testdata/data-dir&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"title"},
				SortAttr: []string{"tags"},
				Data:     `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
			},
			Client: &Client{
				Index: "default",
				DB:    `file://../testdata/data-dir`,
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&db=file://../testdata/data-dir&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"title"},
				SortAttr: []string{"tags"},
				Data:     `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&db=file://../testdata/data-dir&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"*"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&uid=id`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
			},
			Client: &Client{
				Index: "default",
				UID:   "id",
			},
			Search: &Search{},
		},
	},
	cfgTest{
		pt: pt{
			query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&db=file://../testdata/data-dir/audiobooks.json&index=audiobooks&uid=id`,
		},
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				SortAttr:  []string{"tags"},
			},
			Client: &Client{
				UID:   "id",
				Index: "audiobooks",
				DB:    "file://../testdata/data-dir/audiobooks.json",
			},
			Search: &Search{},
		},
	},
}

func TestDecodeCfgStr(t *testing.T) {
	for num, test := range cfgTests {
		cfg := NewCfg()
		err := cfg.Decode(test.query)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SrchAttr", cfg.Idx.SrchAttr, test.Idx.SrchAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "FacetAttr", cfg.Idx.FacetAttr, test.Idx.FacetAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SortAttr", cfg.Idx.SortAttr, test.Idx.SortAttr)
		if err != nil {
			t.Error(err)
		}
		if cfg.IndexName() != test.IndexName() {
			t.Errorf("test %v Index: got %#v, expected %#v\n", num, cfg.IndexName(), test.IndexName())
		}
		if cfg.Client.UID != test.Client.UID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", num, cfg.Client.UID, test.Client.UID)
		}
		if cfg.DataURL().Path != test.DataURL().Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.DataURL().Path, test.DataURL().Path)
		}
		if cfg.DB().Path != test.DB().Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.DB().Path, test.DB().Path)
		}
		if cfg.SrchURL().Path != test.SrchURL().Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.SrchURL().Path, test.SrchURL().Path)
		}
	}
}

func TestDecodeCfgVals(t *testing.T) {
	for num, test := range cfgTests {
		cfg := NewCfg()
		err := cfg.Decode(test.query)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SrchAttr", cfg.Idx.SrchAttr, test.Idx.SrchAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "FacetAttr", cfg.Idx.FacetAttr, test.Idx.FacetAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SortAttr", cfg.Idx.SortAttr, test.Idx.SortAttr)
		if err != nil {
			t.Error(err)
		}
		if cfg.IndexName() != test.IndexName() {
			t.Errorf("test %v Index: got %#v, expected %#v\n", num, cfg.IndexName(), test.IndexName())
		}
		if cfg.Client.UID != test.Client.UID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", num, cfg.Client.UID, test.Client.UID)
		}
		if cfg.DataURL().Path != test.DataURL().Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.DataURL().Path, test.DataURL().Path)
		}
		if cfg.DB().Path != test.DB().Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.DB().Path, test.DB().Path)
		}
		if cfg.SrchURL().Path != test.SrchURL().Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.SrchURL().Path, test.SrchURL().Path)
		}
	}
}

func TestEncodeCfg(t *testing.T) {
	t.SkipNow()
	for num, test := range cfgTests {
		v, err := Encode(test.Idx)
		if err != nil {
			t.Error(err)
		}
		if v.Encode() != test.query {
			t.Errorf("test %v: got %v, wanted %v\n", num, v.Encode(), test.query)
		}
	}
}

func sliceTest(num, field any, got, want []string) error {
	if !slices.Equal(got, want) {
		return paramTestMsg(num, field, got, want)
	}
	return nil
}

func paramTestMsg(num, field, got, want any) error {
	return fmt.Errorf("test %v, field %s\ngot %#v, wanted %#v\n", num, field, got, want)
}
