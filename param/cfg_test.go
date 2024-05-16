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
			SrchAttr: []string{"*"},
			Paramz: &Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=`,
		},
		Cfg: &Cfg{
			SrchAttr: []string{"*"},
			Paramz: &Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title`,
		},
		Cfg: &Cfg{
			SrchAttr: []string{"title"},
			Paramz: &Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&path=file://../testdata/data-dir&sortableAttributes=tags`,
		},
		Cfg: &Cfg{
			SrchAttr: []string{"title"},
			Paramz: &Paramz{
				Index: "default",
				Path:  `file://../testdata/data-dir`,
			},
			SortAttr: []string{"tags"},
		},
	},
	cfgTest{
		pt: pt{
			query: `?attributesForFaceting=tags,authors,series,narrators`,
		},
		Cfg: &Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`,
		},
		Cfg: &Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		},
		Cfg: &Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &Paramz{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&id=id`,
		},
		Cfg: &Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Paramz: &Paramz{
				Index: "default",
				ID:    "id",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&path=file://../testdata/data-dir/audiobooks.json&index=audiobooks&id=id`,
		},
		Cfg: &Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			SortAttr:  []string{"tags"},
			Paramz: &Paramz{
				ID:    "id",
				Index: "audiobooks",
				Path:  "file://../testdata/data-dir/audiobooks.json",
			},
		},
	},
}

func TestDecodeCfgStr(t *testing.T) {
	for num, test := range cfgTests {
		cfg := NewCfg()
		err := Decode(test.query, cfg)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SrchAttr", cfg.SrchAttr, test.SrchAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "FacetAttr", cfg.FacetAttr, test.FacetAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SortAttr", cfg.SortAttr, test.SortAttr)
		if err != nil {
			t.Error(err)
		}
		if cfg.Index != test.Index {
			t.Errorf("test %v Index: got %#v, expected %#v\n", num, cfg.Index, test.Index)
		}
		if cfg.ID != test.ID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", num, cfg.ID, test.ID)
		}
		if cfg.Path != test.Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.Path, test.Path)
		}
	}
}

func TestDecodeCfgVals(t *testing.T) {
	for num, test := range cfgTests {
		cfg := NewCfg()
		err := Decode(test.vals(), cfg)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SrchAttr", cfg.SrchAttr, test.SrchAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "FacetAttr", cfg.FacetAttr, test.FacetAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "SortAttr", cfg.SortAttr, test.SortAttr)
		if err != nil {
			t.Error(err)
		}
		if cfg.Index != test.Index {
			t.Errorf("test %v Index: got %#v, expected %#v\n", num, cfg.Index, test.Index)
		}
		if cfg.ID != test.ID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", num, cfg.ID, test.ID)
		}
		if cfg.Path != test.Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.Path, test.Path)
		}
	}
}

func TestEncodeCfg(t *testing.T) {
	t.SkipNow()
	for num, test := range cfgTests {
		v, err := Encode(test.Cfg)
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
