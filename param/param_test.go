package param

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

type paramTest struct {
	query string
	want  *Params
}

var paramTests = []paramTest{
	paramTest{
		query: ``,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"*"},
			FacetAttr:            []string{""},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "",
			DataFile:             "",
			UID:                  "",
		},
	},
	paramTest{
		query: `?searchableAttributes=`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"*"},
			FacetAttr:            []string{""},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "",
			DataFile:             "",
			UID:                  "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&fullText=../testdata/poot.bleve`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"title"},
			FacetAttr:            []string{""},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "../testdata/poot.bleve",
			DataDir:              "",
			DataFile:             "",
			UID:                  "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&dataDir=../testdata/data-dir`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"title"},
			FacetAttr:            []string{""},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "../testdata/data-dir",
			DataFile:             "",
			UID:                  "",
		},
	},
	paramTest{
		query: `?attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"*"},
			FacetAttr:            []string{"tags", "authors", "series", "narrators"},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "",
			DataFile:             "",
			UID:                  "",
		},
	},
	paramTest{
		query: `?attributesForFaceting=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"*"},
			FacetAttr:            []string{"tags", "authors", "series", "narrators"},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "",
			DataFile:             "../testdata/data-dir/audiobooks.json",
			UID:                  "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"title"},
			FacetAttr:            []string{"tags", "authors", "series", "narrators"},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "",
			DataFile:             "",
			UID:                  "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"title"},
			FacetAttr:            []string{"tags", "authors", "series", "narrators"},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 0,
			HitsPerPage:          0,
			Query:                "",
			SortBy:               "",
			Order:                "",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "",
			DataFile:             "../testdata/data-dir/audiobooks.json",
			UID:                  "",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
		want: &Params{
			Other:                url.Values{},
			SrchAttr:             []string{"title"},
			FacetAttr:            []string{"tags", "authors"},
			SortAttr:             []string{""},
			DefaultField:         "",
			Hits:                 0,
			AttributesToRetrieve: []string{""},
			Page:                 3,
			HitsPerPage:          0,
			Query:                "fish",
			SortBy:               "title",
			Order:                "desc",
			Facets:               []string{""},
			Filters:              "",
			FacetFilters:         []any{""},
			SortFacetsBy:         "",
			BlvPath:              "",
			DataDir:              "",
			DataFile:             "../testdata/data-dir/audiobooks.json",
			UID:                  "",
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
	for _, u := range testQuerySettings {
		p, err := Parse(u)
		if err != nil {
			t.Error(err)
		}
		if !p.IsFile() {
			t.Errorf("url %s has no data file", u)
		}
	}
}

type pathMatch struct {
	prefix string
	path   string
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
	`/blv/../testdata/poot.bleve`: pathMatch{
		prefix: "blv",
		path:   "/home/mxb/code/srch/testdata/poot.bleve",
	},
	`blv/../testdata/poot.bleve`: pathMatch{
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
	`/dir/../testdata/nddata`: pathMatch{
		prefix: "dir",
		path:   "/home/mxb/code/srch/testdata/nddata",
	},
	`dir/../testdata/nddata`: pathMatch{
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
	`/file/../testdata/nddata/ndbooks.ndjson`: pathMatch{
		prefix: "file",
		path:   "/home/mxb/code/srch/testdata/nddata/ndbooks.ndjson",
	},
	`file/../testdata/nddata/ndbooks.ndjson`: pathMatch{
		prefix: "file",
		path:   "/home/mxb/code/srch/testdata/nddata/ndbooks.ndjson",
	},
}

func TestPaths(t *testing.T) {
	for path, want := range pathMatches {
		pre, loc := parsePath(path)
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
		fmt.Printf("\t\tAttributesToRetrieve: %#v,\n", p.AttributesToRetrieve)
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
