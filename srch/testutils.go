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

type CfgTest struct {
	*Cfg
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

func (ct CfgTest) Slice(q string, got, want []string) error {
	if !slices.Equal(got, want) {
		return ct.Err(ct.Msg(q, got, want), errors.New("slices not equal"))
	}
	return nil
}

func (ct CfgTest) Msg(q string, got, want any) string {
	return fmt.Sprintf("%v\ngot %#v, wanted %#v\n", q, got, want)
}

func (ct CfgTest) Err(msg string, err error) error {
	return fmt.Errorf("%v\nerror: %w\n", msg, err)
}

func (ct CfgTest) Srch(got, want *Search) error {
	err := ct.Slice("search.RtrvAttr", got.RtrvAttr, want.RtrvAttr)
	if err != nil {
		return err
	}
	err = ct.Slice("search.Facets", got.Facets, want.Facets)
	if err != nil {
		return err
	}
	err = ct.Slice("search.FacetFltr", got.FacetFltr, want.FacetFltr)
	if err != nil {
		return err
	}
	return nil
}

func (ct CfgTest) Index(got, want *IdxCfg) error {
	err := ct.Slice("search.SrchAttr", got.SrchAttr, want.SrchAttr)
	if err != nil {
		return err
	}
	err = ct.Slice("search.FacetAttr", got.FacetAttr, want.FacetAttr)
	if err != nil {
		return err
	}
	err = ct.Slice("search.SortAttr", got.SortAttr, want.SortAttr)
	if err != nil {
		return err
	}
	return nil
}

func (ct CfgTest) Config(got, want *Cfg) error {
	if got.IndexName() != want.IndexName() {
		return ct.Err(ct.Msg("cfg.IndexName()", got.IndexName(), want.IndexName()), errors.New("index name doesn't match"))
	}
	if got.Client.UID != want.Client.UID {
		return ct.Err(ct.Msg("cfg.Client.UID", got.Client.UID, want.Client.UID), errors.New("index uid doesn't match"))
	}
	if got.DataURL().Path != want.DataURL().Path {
		return ct.Err(ct.Msg("cfg.DataURL().Path", got.DataURL().Path, want.DataURL().Path), errors.New("data path doesn't match"))
	}
	if got.DB().Path != want.DB().Path {
		return ct.Err(ct.Msg("cfg.DB().Path", got.DB().Path, want.DB().Path), errors.New("db path doesn't match"))
	}
	if got.SrchURL().Path != want.SrchURL().Path {
		return ct.Err(ct.Msg("cfg.SrchURL().Path", got.SrchURL().Path, want.SrchURL().Path), errors.New("srch path doesn't match"))
	}
	return nil
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

func ParamTests() map[QueryStr]CfgTest {
	return map[QueryStr]CfgTest{
		QueryStr(``): CfgTest{
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
		QueryStr(`?searchableAttributes=`): CfgTest{
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
		QueryStr(`?searchableAttributes=title`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=series&attributesForFaceting=narrators`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane"]`): CfgTest{
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
		QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&query=fish&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&url=file://home/mxb/code/srch/testdata/hare/audiobooks.json`): CfgTest{
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

func sliceTest(num, field any, got, want []string) error {
	if !slices.Equal(got, want) {
		return paramTestMsg(num, field, got, want)
	}
	return nil
}

func paramTestMsg(num, field, got, want any) error {
	return fmt.Errorf("test %v, field %s\ngot %#v, wanted %#v\n", num, field, got, want)
}
