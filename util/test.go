package util

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ohzqq/srch/param"
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

type ParamCfgTest struct {
	query string
	*param.Cfg
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

var ParamCfgTests = map[QueryStr]ParamCfgTest{
	QueryStr(``): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr: []string{"*"},
			},
			Client: &param.Client{
				Index: "default",
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr: []string{"*"},
			},
			Client: &param.Client{
				Index: "default",
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr: []string{"title"},
			},
			Client: &param.Client{
				Index: "default",
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr: []string{"title"},
				SortAttr: []string{"tags"},
				Data:     DataTestURL,
			},
			Client: &param.Client{
				Index: "default",
				DB:    HareTestURL,
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "default",
				DB:    HareTestURL,
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "default",
				DB:    HareTestURL,
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "default",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &param.Search{
				RtrvAttr: []string{"*"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &param.Search{
				RtrvAttr: []string{"title", "tags", "authors"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &param.Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &param.Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
				Query:    "fish",
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane"]`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &param.Search{
				RtrvAttr:  []string{"title", "tags", "authors"},
				Facets:    []string{"tags", "authors", "series", "narrators"},
				Query:     "fish",
				FacetFltr: []string{"[\"authors:amy lane\"]"},
			},
		},
	},
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&url=file://home/mxb/code/srch/testdata/hare/audiobooks.json`): ParamCfgTest{
		Cfg: &param.Cfg{
			Idx: &param.Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"tags"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Data:      DataTestURL,
			},
			Client: &param.Client{
				Index: "audiobooks",
				DB:    HareTestURL,
				UID:   "id",
			},
			Search: &param.Search{
				RtrvAttr:  []string{"title", "tags", "authors"},
				Facets:    []string{"tags", "authors", "series", "narrators"},
				Query:     "fish",
				FacetFltr: []string{"[\"authors:amy lane\", [\"tags:romance\", \"tags:-dnr\"]]"},
				URI:       filepath.Join(HareTestURL, "audiobooks.json"),
			},
		},
	},
}
