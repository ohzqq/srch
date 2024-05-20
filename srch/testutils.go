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

func newTestReq(v any) (reqTest, error) {
	req, err := NewRequest(v)
	if err != nil {
		return reqTest{}, err
	}
	return reqTest{Request: req}, nil
}

func (t reqTest) cfgTest(idx int) cfgTest {
	cfg := getTestCfg(idx)
	return cfgTest{Cfg: cfg}
}

func (t reqTest) clientTest(idx int) clientTest {
	client := getTestClient(idx)
	return clientTest{Client: client}
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

func getTestCfg(idx int) *Cfg {
	tests := []*Cfg{
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
		&Cfg{
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
	}
	return tests[idx]
}
