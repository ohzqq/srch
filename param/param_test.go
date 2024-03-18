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
		query: `searchableAttributes=`,
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
		query: `searchableAttributes=title&fullText=../testdata/poot.bleve`,
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
		query: `searchableAttributes=title&dataDir=../testdata/data-dir`,
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
		query: `attributesForFaceting=tags,authors,series,narrators`,
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
		query: `attributesForFaceting=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json`,
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
		query: `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
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
		query: `searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators`,
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
		query: `searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
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
