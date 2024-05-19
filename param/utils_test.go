package param

import (
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

const (
	HareTestPath  = `/home/mxb/code/srch/testdata/hare`
	HareTestURL   = `file://home/mxb/code/srch/testdata/hare`
	HareTestQuery = `?url=file://home/mxb/code/srch/testdata/hare/`
)

const (
	DataTestURL = `file://home/mxb/code/srch/testdata/ndbooks.ndjson`
	IdxTestFile = `file://home/mxb/code/srch/testdata/hare/audiobooks.json`
)

type QueryStr string

type CfgTest struct {
	*Cfg
}

func (p QueryStr) String() string {
	return string(p)
}

func (p QueryStr) Query() url.Values {
	v, _ := url.ParseQuery(strings.TrimPrefix(p.String(), "?"))
	return v
}

func (p QueryStr) URL() *url.URL {
	u, _ := url.Parse(p.String())
	return u
}

func SrchTests(t *testing.T, num QueryStr, got, want *Search) {
	err := sliceTest(num, "RtrvAttr", got.RtrvAttr, want.RtrvAttr)
	if err != nil {
		t.Error(err)
	}
}

func CfgTests(t *testing.T, num QueryStr, got, want *Cfg) {
	if got.IndexName() != want.IndexName() {
		t.Errorf("test %v Index: got %#v, expected %#v\n", num, got.IndexName(), want.IndexName())
	}
	if got.Client.UID != want.Client.UID {
		t.Errorf("test %v ID: got %#v, expected %#v\n", num, got.Client.UID, want.Client.UID)
	}
	if got.DataURL().Path != want.DataURL().Path {
		t.Errorf("test %v Path: got %#v, expected %#v\n", num, got.DataURL().Path, want.DataURL().Path)
	}
	if got.DB().Path != want.DB().Path {
		t.Errorf("test %v Path: got %#v, expected %#v\n", num, got.DB().Path, want.DB().Path)
	}
	if got.SrchURL().Path != want.SrchURL().Path {
		t.Errorf("test %v Path: got %#v, expected %#v\n", num, got.SrchURL().Path, want.SrchURL().Path)
	}
}

func IdxTests(t *testing.T, num QueryStr, got, want *Idx) {
	err := sliceTest(num, "SrchAttr", got.SrchAttr, want.SrchAttr)
	if err != nil {
		t.Error(err)
	}
	err = sliceTest(num, "FacetAttr", got.FacetAttr, want.FacetAttr)
	if err != nil {
		t.Error(err)
	}
	err = sliceTest(num, "SortAttr", got.SortAttr, want.SortAttr)
	if err != nil {
		t.Error(err)
	}
}

var TestQueryParams = []QueryStr{
	QueryStr(``),
	QueryStr(`?searchableAttributes=`),
	QueryStr(`?searchableAttributes=title`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane"]`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&url=file://home/mxb/code/srch/testdata/hare/audiobooks.json`),
}

var cfgTests = map[QueryStr]CfgTest{
	QueryStr(``): CfgTest{
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
	QueryStr(`?searchableAttributes=`): CfgTest{
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
	QueryStr(`?searchableAttributes=title`): CfgTest{
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
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr: []string{"title"},
				SortAttr: []string{"tags"},
				Data:     DataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    HareTestURL,
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    HareTestURL,
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    HareTestURL,
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "default",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
				Query:    "fish",
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane"]`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    HareTestURL,
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
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&url=file://home/mxb/code/srch/testdata/hare/audiobooks.json`): CfgTest{
		Cfg: &Cfg{
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &Search{
				RtrvAttr:  []string{"title", "tags", "authors"},
				Facets:    []string{"tags", "authors", "series", "narrators"},
				Query:     "fish",
				FacetFltr: []string{"[\"authors:amy lane\", [\"tags:romance\", \"tags:-dnr\"]]"},
				URI:       filepath.Join(HareTestURL, "audiobooks.json"),
			},
		},
	},
}
