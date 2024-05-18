package param

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"
)

var cfgTests = []cfgTest{
	cfgTest{
		query: ``,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"*"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"*"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"title"},
			},
			Client: &Client{
				Index: "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"title"},
				SortAttr: []string{"tags"},
				Data:     dataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    hareTestURL,
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    hareTestURL,
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    hareTestURL,
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    hareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    hareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    hareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    hareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    hareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
				Query:    "fish",
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane"]`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    hareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr:  []string{"title", "tags", "authors"},
				Facets:    []string{"tags", "authors", "series", "narrators"},
				Query:     "fish",
				FacetFltr: []string{"[\"authors:amy lane\"]"},
			},
		},
	},
	cfgTest{
		query: `?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&url=file://home/mxb/code/srch/testdata/hare/audiobooks.json`,
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      dataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    hareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr:  []string{"title", "tags", "authors"},
				Facets:    []string{"tags", "authors", "series", "narrators"},
				Query:     "fish",
				FacetFltr: []string{"[\"authors:amy lane\", [\"tags:romance\", \"tags:-dnr\"]]"},
				URI:       filepath.Join(hareTestURL, "audiobooks.json"),
			},
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

		testIdx(t, num, cfg.Idx, test.Idx)
		testCfg(t, num, cfg, test.Cfg)
		testSrch(t, num, cfg.Search, test.Search)

	}
}

func TestDecodeCfgVals(t *testing.T) {
	for num, test := range cfgTests {
		cfg := NewCfg()
		err := cfg.Decode(test.query)
		if err != nil {
			t.Error(err)
		}

		testIdx(t, num, cfg.Idx, test.Idx)

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
