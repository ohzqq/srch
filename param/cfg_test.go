package param

import (
	"fmt"
	"slices"
	"testing"
)

type cfgTest struct {
	query string
	*Cfg
}

var cfgTests = []cfgTest{
	cfgTest{
		query: ``,
		Cfg: &Cfg{
			SrchAttr: []string{"*"},
			Index:    "default",
		},
	},
	cfgTest{
		query: `?searchableAttributes=`,
		Cfg: &Cfg{
			SrchAttr: []string{"*"},
			Index:    "default",
		},
	},
	cfgTest{
		query: `?searchableAttributes=title`,
		Cfg: &Cfg{
			SrchAttr: []string{"title"},
			Index:    "default",
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&path=../testdata/data-dir&sortableAttributes=tags`,
		Cfg: &Cfg{
			SrchAttr: []string{"title"},
			Index:    "default",
			Path:     `../testdata/data-dir`,
			SortAttr: []string{"tags"},
		},
	},
	cfgTest{
		query: `?attributesForFaceting=tags,authors,series,narrators`,
		Cfg: &Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Index:     "default",
		},
	},
	cfgTest{
		query: `?attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`,
		Cfg: &Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Index:     "default",
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		Cfg: &Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Index:     "default",
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&id=id`,
		Cfg: &Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Index:     "default",
			ID:        "id",
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&path=../testdata/data-dir/audiobooks.json&index=audiobooks&id=id`,
		Cfg: &Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			SortAttr:  []string{"tags"},
			ID:        "id",
			Index:     "audiobooks",
			Path:      "../testdata/data-dir/audiobooks.json",
		},
	},
}

func TestDecodeCfg(t *testing.T) {
	for num, test := range cfgTests {
		cfg, err := ParseCfg(test.query)
		if err != nil {
			t.Error(err)
		}
		if !slices.Equal(cfg.SrchAttr, test.SrchAttr) {
			t.Errorf("test %v SrchAttr: got %#v, expected %#v\n", num, cfg.SrchAttr, test.SrchAttr)
		}
		if !slices.Equal(cfg.FacetAttr, test.FacetAttr) {
			t.Errorf("test %v FacetAttr: got %#v, expected %#v\n", num, cfg.FacetAttr, test.FacetAttr)
		}
		if !slices.Equal(cfg.SortAttr, test.SortAttr) {
			t.Errorf("test %v SortAttr: got %#v, expected %#v\n", num, cfg.SortAttr, test.SortAttr)
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
