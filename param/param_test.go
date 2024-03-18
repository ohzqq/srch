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
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"*"}, FacetAttr: []string{""}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "", DataFile: []string{""}, UID: ""},
		},
	},
	paramTest{
		query: `searchableAttributes=`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"*"}, FacetAttr: []string{""}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "", DataFile: []string{""}, UID: ""},
		},
	},
	paramTest{
		query: `searchableAttributes=title&fullText=../testdata/poot.bleve`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"title"}, FacetAttr: []string{""}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "../testdata/poot.bleve", DataDir: "", DataFile: []string{""}, UID: ""},
		},
	},
	paramTest{
		query: `searchableAttributes=title&dataDir=../testdata/data-dir`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"title"}, FacetAttr: []string{""}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "../testdata/data-dir", DataFile: []string{""}, UID: ""},
		},
	},
	paramTest{
		query: `attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"*"}, FacetAttr: []string{"tags", "authors", "series", "narrators"}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "", DataFile: []string{""}, UID: ""},
		},
	},
	paramTest{
		query: `attributesForFaceting=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"*"}, FacetAttr: []string{"tags", "authors", "series", "narrators"}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "", DataFile: []string{"../testdata/data-dir/audiobooks.json"}, UID: ""},
		},
	},
	paramTest{
		query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"title"}, FacetAttr: []string{"tags", "authors", "series", "narrators"}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "", DataFile: []string{""}, UID: ""},
		},
	},
	paramTest{
		query: `searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"title"}, FacetAttr: []string{"tags", "authors", "series", "narrators"}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 0,
				HitsPerPage:          0,
				Query:                "",
				SortBy:               "",
				Order:                "",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "", DataFile: []string{"../testdata/data-dir/audiobooks.json"}, UID: ""},
		},
	},
	paramTest{
		query: `searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
		want: &Params{
			Other:         url.Values{},
			IndexSettings: &IndexSettings{SrchAttr: []string{"title"}, FacetAttr: []string{"tags", "authors"}, SortAttr: []string{""}, DefaultField: "", UID: ""},
			Search: &Search{
				Hits:                 0,
				AttributesToRetrieve: []string{""},
				Page:                 3,
				HitsPerPage:          0,
				Query:                "fish",
				SortBy:               "title",
				Order:                "desc",
				FacetSettings:        &FacetSettings{UID: "", Facets: []string{""}, Filters: "", FacetFilters: []any{""}, SortFacetsBy: ""},
			},
			SrchCfg: &SrchCfg{BlvPath: "", DataDir: "", DataFile: []string{"../testdata/data-dir/audiobooks.json"}, UID: ""},
		},
	},
}

var testQuerySettings = []string{
	"",
	"searchableAttributes=",
	"searchableAttributes=title&fullText=../testdata/poot.bleve",
	"searchableAttributes=title&dataDir=../testdata/data-dir",
	"attributesForFaceting=tags,authors,series,narrators",
	"attributesForFaceting=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators",
	"searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators",
	`searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
}

func TestNewRequest(t *testing.T) {
	for _, test := range testQuerySettings {
		req, err := NewClient(test)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%#v\n", req)
	}
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
				t.Errorf("got %#v, exptect %#v\n", sa, attr)
			}
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
		fmt.Printf("\tIndexSettings: %#v,\n", p.IndexSettings)
		fmt.Println("\tSearch: &Search{")
		fmt.Printf("\t\tHits: %#v,\n", p.Search.Hits)
		fmt.Printf("\t\tAttributesToRetrieve: %#v,\n", p.Search.AttributesToRetrieve)
		fmt.Printf("\t\tPage: %#v,\n", p.Search.Page)
		fmt.Printf("\t\tHitsPerPage: %#v,\n", p.Search.HitsPerPage)
		fmt.Printf("\t\tQuery: %#v,\n", p.Search.Query)
		fmt.Printf("\t\tSortBy: %#v,\n", p.Search.SortBy)
		fmt.Printf("\t\tOrder: %#v,\n", p.Search.Order)
		fmt.Printf("\t\tFacetSettings: %#v,\n", p.Search.FacetSettings)
		println("},")
		fmt.Printf("\tSrchCfg: %#v,\n", p.SrchCfg)
		println("},")
		println("},")
	}
}
