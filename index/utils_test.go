package index

import (
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"github.com/ohzqq/srch/param"
)

const hareTestPath = `/home/mxb/code/srch/testdata/hare`
const hareTestURL = `file://home/mxb/code/srch/testdata/hare`
const hareTestQuery = `?url=file://home/mxb/code/srch/testdata/hare/`

var defTbls = []string{"_settings"}

var cfgTests = []test{
	test{
		query: ``,
		Cfg: &param.Cfg{
			SrchAttr: []string{"*"},
			Client: &param.Client{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title`,
		Cfg: &param.Cfg{
			SrchAttr: []string{"title"},
			Client: &param.Client{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title&url=file://home/mxb/code/srch/testdata/hare/&sortableAttributes=tags`,
		Cfg: &param.Cfg{
			SrchAttr: []string{"title"},
			Client: &param.Client{
				Index: "default",
				DB:    `file://home/mxb/code/srch/testdata/hare/`,
			},
			SortAttr: []string{"tags"},
		},
	},
	test{
		query: `?attributesForFaceting=tags,authors,series,narrators`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"*"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Client: &param.Client{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Client: &param.Client{
				Index: "default",
			},
		},
	},
	test{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&uid=id&index=audiobooks`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			Client: &param.Client{
				Index: "audiobooks",
				UID:   "id",
			},
		},
	},
	test{
		query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&db=/home/mxb/code/srch/testdata/hare/&uid=id&index=audiobooks`,
		//query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&sortableAttributes=tags&url=file://home/mxb/code/srch/testdata/hare/&uid=id`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			SortAttr:  []string{"tags"},
			Client: &param.Client{
				UID:   "id",
				Index: "audiobooks",
				DB:    "/home/mxb/code/srch/testdata/hare/",
			},
		},
	},
	test{
		query: `searchableAttributes=title&attributesForFaceting=tags,authors,series&sortableAttributes=tags&db=/home/mxb/code/srch/testdata/hare/&uid=id&index=audiobooks`,
		Cfg: &param.Cfg{
			SrchAttr:  []string{"title"},
			FacetAttr: []string{"tags", "authors", "series", "narrators"},
			SortAttr:  []string{"tags"},
			Client: &param.Client{
				UID:   "id",
				Index: "audiobooks",
				DB:    "/home/mxb/code/srch/testdata/hare/",
			},
		},
	},
}

type params struct {
	query string
	Cfg   *param.Cfg
}

type test struct {
	query string
	Cfg   *param.Cfg
}

func (p test) str() string {
	return p.query
}

func (p test) vals() url.Values {
	v, _ := url.ParseQuery(strings.TrimPrefix(p.query, "?"))
	return v
}

func (p test) url() *url.URL {
	u, err := url.Parse(url.QueryEscape(p.query))
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func (p test) slice(got, want []string, err error) error {
	if !slices.Equal(got, want) {
		return p.err(got, want, err)
	}
	return nil
}

func (p test) err(got, want any, err error) error {
	return fmt.Errorf("query %v\ngot %#v, wanted %#v\nerror: %w\n", p.str(), got, want, err)
}

func (t test) msg(msg any) error {
	return fmt.Errorf("query: %v\nerror: %#v\n", t.str(), msg)
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
