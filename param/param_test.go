package param

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type paramTest struct {
	query string
	want  *Params
}

var paramTests = []paramTest{
	paramTest{
		query: ``,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?searchableAttributes=`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&fullText=../testdata/poot.bleve`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&dataDir=../testdata/data-dir`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?attributesForFaceting=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
		want: &Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{"tags", "authors"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         3,
			HitsPerPage:  0,
			Query:        "fish",
			SortBy:       "title",
			Order:        "desc",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
		},
	},
}

var testQuerySettings = []string{
	"blv/../testdata/poot.bleve?searchableAttributes=title&facets=tags,authors,series,narrators",
	"/dir/home/mxb/code/srch/testdata/data-dir?searchableAttributes=title&facets=tags,authors,series,narrators",
	"file/home/mxb/code/srch/testdata/data-dir/audiobooks.json?searchableAttributes=title&facets=tags,authors,series,narrators",
	`/file/home/mxb/code/srch/testdata/data-dir/audiobooks.json?searchableAttributes=title&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
}

func TestNewParams(t *testing.T) {
	t.SkipNow()
	for i, test := range paramTests {
		p, err := Parse(test.query)
		if err != nil {
			t.Error(err)
		}
		//println(test.query)
		if i > 1 {
			attr := p.SrchAttr[0]
			if sa := test.want.SrchAttr[0]; sa != attr {
				t.Errorf("%d: test query %s\ngot %#v, exptect %#v\n", i, test.query, attr, sa)
			}
		}
	}
}

func TestNewQueryURLs(t *testing.T) {
	t.SkipNow()
	for _, u := range testQuerySettings {
		p, err := Parse(u)
		if err != nil {
			t.Error(err)
		}
		println(p.String())
	}
}

func TestDecodeParams(t *testing.T) {
	t.SkipNow()
	for _, key := range SettingParams {
		switch key {
		case SrchAttr:
			viper.SetDefault(key.Snake(), []string{"title"})
		case FacetAttr:
			viper.SetDefault(key.Snake(), []string{"tags"})
		case SortAttr:
			viper.SetDefault(key.Snake(), []string{"title:desc"})
		case UID:
			viper.SetDefault(key.Snake(), "id")
		}
	}

	for _, key := range SearchParams {
		switch key {
		case SortFacetsBy:
			viper.SetDefault(key.Snake(), "tags:count:desc")
		case Facets:
			viper.SetDefault(key.Snake(), []string{"tags"})
		case RtrvAttr:
			viper.SetDefault(key.Snake(), []string{"*"})
		case Page:
			viper.SetDefault(key.Snake(), 0)
		case HitsPerPage:
			viper.SetDefault(key.Snake(), -1)
		case SortBy:
			viper.SetDefault(key.Snake(), "title")
		case Order:
			viper.SetDefault(key.Snake(), "desc")
		}
	}

	s := viper.AllSettings()

	dec, err := DecodeSnakeMap(s)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", dec)
}

func TestParamStringer(t *testing.T) {
	t.SkipNow()
	test := paramTests[len(paramTests)-1]
	p, err := Parse(test.query)
	if err != nil {
		t.Error(err)
	}
	s := make(map[string]any)
	for k, v := range p.Values() {
		s[k] = v
	}
	q := QueryToSettings(s)

	search := lo.Map(SearchParams, func(i Param, _ int) string { return i.Snake() })
	settings := lo.Map(SettingParams, func(i Param, _ int) string { return i.Snake() })
	keys := maps.Keys(q)

	ss := lo.Intersect(keys, search)
	if len(ss) < 1 {
		t.Errorf("not same %#v\n", ss)
	}
	ss = lo.Intersect(keys, settings)
	if len(ss) < 1 {
		t.Errorf("not same %#v\n", ss)
	}

}

type pathMatch struct {
	prefix      string
	path        string
	contentType string
}

var pathMatches = map[string]pathMatch{
	``: pathMatch{
		prefix: "",
		path:   "",
	},
	`/`: pathMatch{
		prefix: "",
		path:   "",
	},
	`/blv`: pathMatch{
		prefix: "blv",
		path:   "",
	},
	`blv`: pathMatch{
		prefix: "blv",
		path:   "",
	},
	`/blv?path=../testdata/poot.bleve`: pathMatch{
		prefix: "blv",
		path:   "/home/mxb/code/srch/testdata/poot.bleve",
	},
	`blv?path=../testdata/poot.bleve`: pathMatch{
		prefix: "blv",
		path:   "/home/mxb/code/srch/testdata/poot.bleve",
	},
	`/dir`: pathMatch{
		prefix: "dir",
		path:   "",
	},
	`dir`: pathMatch{
		prefix: "dir",
		path:   "",
	},
	`/dir?path=../testdata/nddata`: pathMatch{
		prefix: "dir",
		path:   "/home/mxb/code/srch/testdata/nddata",
	},
	`dir?path=../testdata/nddata`: pathMatch{
		prefix: "dir",
		path:   "/home/mxb/code/srch/testdata/nddata",
	},
	`/file`: pathMatch{
		prefix: "file",
		path:   "",
	},
	`file`: pathMatch{
		prefix: "file",
		path:   "",
	},
	`/file?path=../testdata/nddata/ndbooks.ndjson`: pathMatch{
		prefix:      "file",
		path:        "/home/mxb/code/srch/testdata/nddata/ndbooks.ndjson",
		contentType: NdJSON,
	},
	`file?path=../testdata/nddata/ndbooks`: pathMatch{
		prefix: "file",
		path:   "/home/mxb/code/srch/testdata/nddata/ndbooks",
	},
	`file?path=../testdata/data-dir/audiobooks.json`: pathMatch{
		prefix:      "file",
		path:        "/home/mxb/code/srch/testdata/data-dir/audiobooks.json",
		contentType: JSON,
	},
	`blv?path=/mnt/x/libraries/audiobooks/audiobooks.bleve`: pathMatch{
		prefix: "blv",
		path:   "/mnt/x/libraries/audiobooks/audiobooks.bleve",
	},
}

func TestAbsPath(t *testing.T) {
	t.SkipNow()
	path := `blv/mnt/x/libraries/audiobooks/audiobooks.bleve`
	if !filepath.IsAbs(path) {
		println("not abs")
	}
	matches := pathRegexp.FindStringSubmatch(path)
	ri := pathRegexp.SubexpIndex("route")
	pi := pathRegexp.SubexpIndex("path")
	println(matches[ri])
	println(matches[pi])
}

func TestPaths(t *testing.T) {
	t.SkipNow()
	for path, want := range pathMatches {
		params, err := Parse(path)
		if err != nil {
			t.Error(err)
		}
		pre, loc := params.Route, params.Path
		if loc != "" && (want.prefix != pre || loc != want.path) {
			t.Errorf("pre %s, path %s: wnat %#v", pre, loc, want)
		}
	}
}

func printTests() {
	for _, test := range testQuerySettings {
		p, err := Parse(test)
		if err != nil {
			log.Fatal(err)
		}

		println("paramTest{")
		fmt.Printf("test: `%s`,\n", test)
		println("want: &Params{")
		fmt.Println("\tSearch: &Search{")
		fmt.Printf("\t\tHits: %#v,\n", p.Hits)
		fmt.Printf("\t\tPage: %#v,\n", p.Page)
		fmt.Printf("\t\tHitsPerPage: %#v,\n", p.HitsPerPage)
		fmt.Printf("\t\tQuery: %#v,\n", p.Query)
		fmt.Printf("\t\tSortBy: %#v,\n", p.SortBy)
		fmt.Printf("\t\tOrder: %#v,\n", p.Order)
		println("},")
		println("},")
		println("},")
	}
}
