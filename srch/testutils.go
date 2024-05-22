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
	HareTestQuery = `?idx=file://home/mxb/code/srch/testdata/hare/`
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

func (t reqTest) getCfg(idx int) *Cfg {
	return getTestCfg(idx)
}

func (t reqTest) cfgWant(idx int) cfgTest {
	return cfgTest{Cfg: t.getCfg(idx)}
}

func (t reqTest) cfgGot() (*Cfg, error) {
	return t.Cfg()
}

func (t reqTest) getClientWant(idx int) *Client {
	client, _ := NewClient(t.getCfg(idx))
	return client
}

func (t reqTest) clientTest(idx int) (clientTest, error) {
	g, err := t.Client()
	if err != nil {
		return clientTest{}, err
	}
	return clientTest{
		got:  g,
		want: t.getClientWant(idx),
	}, nil
}

func (t reqTest) clientWant(idx int) clientTest {
	return clientTest{Client: t.getClientWant(idx)}
}

func (t reqTest) clientGot() (*Client, error) {
	return t.Client()
}

func (t reqTest) getTestIdx(q QueryStr) int {
	idx := slices.Index(TestQueryParams, q)
	if idx == -1 {
		return 0
	}
	return idx
}

func (t reqTest) getQuery(idx int) QueryStr {
	for i := range TestQueryParams {
		if i == idx {
			return TestQueryParams[idx]
		}
	}
	return TestQueryParams[0]
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

func strSliceErr(q string, got, want []string) error {
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		return newErr(msg(q, got, want), errors.New("slices not equal"))
	}
	return nil
}

func intSliceErr(q string, got, want []int) error {
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		return newErr(msg(q, got, want), errors.New("slices not equal"))
	}
	return nil
}

func msg(q string, got, want any) string {
	return fmt.Sprintf("%v\ngot %#v, wanted %#v\n", q, got, want)
}

func newErr(msg string, err error) error {
	return fmt.Errorf("%v\nerror: %w\n", msg, err)
}

var TestQueryParams = []QueryStr{
	QueryStr(``),
	QueryStr(`?searchableAttributes=`),
	QueryStr(`?searchableAttributes=title`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&name=audiobooks`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&name=audiobooks&attributesToRetrieve=title,tags,authors`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&name=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&name=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&name=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane"]`),
	QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&name=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&idx=file://home/mxb/code/srch/testdata/hare/audiobooks.json`),
}

func getTestCfg(idx int) *Cfg {
	tests := []*Cfg{
		&Cfg{
			Idx: &Idx{
				SrchAttr: []string{"*"},
				Name:     "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Idx: &Idx{
				SrchAttr: []string{"*"},
				Name:     "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Idx: &Idx{
				SrchAttr: []string{"title"},
				Name:     "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr: []string{"title"},
				SortAttr: []string{"title"},
				Name:     "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "default",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "default",
				UID:       "id",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "audiobooks",
				UID:       "id",
			},
			Search: &Search{
				RtrvAttr: []string{"*"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "audiobooks",
				UID:       "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "audiobooks",
				UID:       "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "audiobooks",
				UID:       "id",
			},
			Search: &Search{
				RtrvAttr: []string{"title", "tags", "authors"},
				Facets:   []string{"tags", "authors", "series", "narrators"},
				Query:    "fish",
			},
		},
		&Cfg{
			Data: DataTestURL,
			Hare: HareTestURL,
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "audiobooks",
				UID:       "id",
			},
			Search: &Search{
				RtrvAttr:  []string{"title", "tags", "authors"},
				Facets:    []string{"tags", "authors", "series", "narrators"},
				Query:     "fish",
				FacetFltr: []string{"[\"authors:amy lane\"]"},
			},
		},
		&Cfg{
			Data:   DataTestURL,
			Hare:   HareTestURL,
			IdxURL: filepath.Join(HareTestURL, "audiobooks.json"),
			Idx: &Idx{
				SrchAttr:  []string{"title"},
				SortAttr:  []string{"title"},
				FacetAttr: []string{"tags", "authors", "series", "narrators"},
				Name:      "audiobooks",
				UID:       "id",
			},
			Search: &Search{
				RtrvAttr:  []string{"title", "tags", "authors"},
				Facets:    []string{"tags", "authors", "series", "narrators"},
				Query:     "fish",
				FacetFltr: []string{"[\"authors:amy lane\", [\"tags:romance\", \"tags:-dnr\"]]"},
			},
		},
	}
	return tests[idx]
}
