package srch

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
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

type reqTest struct {
	*Request
}

type clientTest struct {
	*Client
}

func newTestReq(v any) (reqTest, error) {
	req, err := NewRequest(v)
	if err != nil {
		return reqTest{}, err
	}
	return reqTest{Request: req}, nil
}

func (t reqTest) cfgTest(want *Cfg) cfgTest {
	return cfgTest{Cfg: want}
}

func (t reqTest) clientTest(want *Client) clientTest {
	return clientTest{Client: want}
}

func (q QueryStr) String() string {
	return string(q)
}

func (q QueryStr) Query() url.Values {
	v, _ := url.ParseQuery(strings.TrimPrefix(q.String(), "?"))
	return v
}

func (q QueryStr) URL() *url.URL {
	u, _ := url.Parse(q.String())
	return u
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

func clientTests() map[QueryStr]cfgTest {
	return map[QueryStr]cfgTest{
		TestQueryParams[0]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr: []string{"*"},
				},
				Client: &ClientCfg{
					Index: "default",
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[1]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr: []string{"*"},
				},
				Client: &ClientCfg{
					Index: "default",
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[2]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr: []string{"title"},
				},
				Client: &ClientCfg{
					Index: "default",
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[3]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr: []string{"title"},
					SortAttr: []string{"tags"},
					Data:     DataTestURL,
				},
				Client: &ClientCfg{
					Index: "default",
					DB:    HareTestURL,
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[4]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
					Index: "default",
					DB:    HareTestURL,
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[5]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
					Index: "default",
					DB:    HareTestURL,
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[6]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
					Index: "default",
					DB:    HareTestURL,
					UID:   "id",
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[7]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
					Index: "audiobooks",
					DB:    HareTestURL,
					UID:   "id",
				},
				Search: &Search{
					RtrvAttr: []string{"*"},
				},
			},
		},
		TestQueryParams[8]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
					Index: "audiobooks",
					DB:    HareTestURL,
					UID:   "id",
				},
				Search: &Search{
					RtrvAttr: []string{"title", "tags", "authors"},
				},
			},
		},
		TestQueryParams[9]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
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
		TestQueryParams[10]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
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
		TestQueryParams[11]: cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&url=file://home/mxb/code/srch/testdata/hare/audiobooks.json`): cfgTest{
			Cfg: &Cfg{
				Idx: &IdxCfg{
					SrchAttr:  []string{"title"},
					SortAttr:  []string{"tags"},
					FacetAttr: []string{"tags", "authors", "series", "narrators"},
					Data:      DataTestURL,
				},
				Client: &ClientCfg{
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
}

func sliceErr(q string, got, want []string) error {
	if !slices.Equal(got, want) {
		return err(msg(q, got, want), errors.New("slices not equal"))
	}
	return nil
}

func msg(q string, got, want any) string {
	return fmt.Sprintf("%v\ngot %#v, wanted %#v\n", q, got, want)
}

func err(msg string, err error) error {
	return fmt.Errorf("%v\nerror: %w\n", msg, err)
}
