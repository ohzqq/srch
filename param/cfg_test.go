package param

import (
	"fmt"
	"slices"
	"testing"
)

type cfgTest struct {
	pt
	*Idx
}

var cfgTests = []cfgTest{
	cfgTest{
		pt: pt{
			query: ``,
		},
		Idx: &Idx{
			SrchAttr: []string{"*"},
			Client: &Client{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=`,
		},
		Idx: &Idx{
			SrchAttr: []string{"*"},
			Client: &Client{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title`,
		},
		Idx: &Idx{
			SrchAttr: []string{"title"},
			Client: &Client{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&url=file://../testdata/data-dir&sortableAttributes=tags`,
		},
		Idx: &Idx{
			SrchAttr: []string{"title"},
			Client: &Client{
				Index: "default",
				DB:    `file://../testdata/data-dir`,
			},
			SortAttr: []string{"tags"},
		},
	},
	cfgTest{
		pt: pt{
			query: `?attributesForFaceting=tags,authors,series,narrators`,
		},
		Idx: &Idx{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Client: &Client{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`,
		},
		Idx: &Idx{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Client: &Client{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		},
		Idx: &Idx{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Client: &Client{
				Index: "default",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&id=id`,
		},
		Idx: &Idx{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Client: &Client{
				Index: "default",
				UID:   "id",
			},
		},
	},
	cfgTest{
		pt: pt{
			query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&url=file://../testdata/data-dir/audiobooks.json&index=audiobooks&id=id`,
		},
		Idx: &Idx{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			SortAttr:  []string{"tags"},
			Client: &Client{
				UID:   "id",
				Index: "audiobooks",
				DB:    "file://../testdata/data-dir/audiobooks.json",
			},
		},
	},
}

func TestDecodeCfgStr(t *testing.T) {
	for num, test := range cfgTests {
		cfg := NewIdx()
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
		if cfg.UID != test.UID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", num, cfg.UID, test.UID)
		}
		if cfg.URI != test.URI {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.URI, test.URI)
		}
	}
}

func TestDecodeCfgVals(t *testing.T) {
	for num, test := range cfgTests {
		cfg := NewIdx()
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
		if cfg.UID != test.UID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", num, cfg.UID, test.UID)
		}
		if cfg.URI != test.URI {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, cfg.URI, test.URI)
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
