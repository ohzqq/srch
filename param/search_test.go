package param

import "testing"

type searchTest struct {
	query string
	*Search
}

var srchTests = []searchTest{
	searchTest{
		query: ``,
		Search: &Search{
			RtrvAttr: []string{"*"},
			Index:    "default",
		},
	},
	searchTest{
		query: `?path=/home/mxb/code/srch/testdata/ndbooks.ndjson`,
		Search: &Search{
			RtrvAttr: []string{"*"},
			Path:     `/home/mxb/code/srch/testdata/ndbooks.ndjson`,
			Index:    "default",
		},
	},
	searchTest{
		query: `?path=/home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish`,
		Search: &Search{
			RtrvAttr: []string{"*"},
			Query:    "fish",
			Path:     `/home/mxb/code/srch/testdata/ndbooks.ndjson`,
			Index:    "default",
		},
	},
	searchTest{
		query: `?path=/home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors`,
		Search: &Search{
			Path:     `/home/mxb/code/srch/testdata/ndbooks.ndjson`,
			Query:    "fish",
			RtrvAttr: []string{"title", "tags", "authors"},
			Index:    "default",
		},
	},
	searchTest{
		query: `?path=/home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title&attributesToRetrieve=tags&attributesToRetrieve=authors`,
		Search: &Search{
			Path:     `/home/mxb/code/srch/testdata/ndbooks.ndjson`,
			Query:    "fish",
			RtrvAttr: []string{"title", "tags", "authors"},
			Index:    "default",
		},
	},
	searchTest{
		query: `?path=/home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`,
		Search: &Search{
			Path:     `/home/mxb/code/srch/testdata/ndbooks.ndjson`,
			Query:    "fish",
			RtrvAttr: []string{"title", "tags", "authors"},
			Facets:   []string{"tags", "authors", "series", "narrators"},
			Index:    "default",
		},
	},
	searchTest{
		query: `?path=/home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&facetFilters=["authors:amy lane"]`,
		Search: &Search{
			Path:      `/home/mxb/code/srch/testdata/ndbooks.ndjson`,
			Query:     "fish",
			RtrvAttr:  []string{"title", "tags", "authors"},
			Facets:    []string{"tags", "authors", "series", "narrators"},
			FacetFltr: []string{"[\"authors:amy lane\"]"},
			Index:     "default",
		},
	},
	searchTest{
		query: `??path=/home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&page=3&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
		Search: &Search{
			RtrvAttr:  []string{"title", "tags", "authors"},
			Page:      3,
			Query:     "fish",
			SortBy:    "title",
			Order:     "desc",
			Facets:    []string{"tags", "authors", "series", "narrators"},
			FacetFltr: []string{"[\"authors:amy lane\", [\"tags:romance\", \"tags:-dnr\"]]"},
			Index:     "default",
		},
	},
}

func TestDecodeSearch(t *testing.T) {
	for num, test := range srchTests {
		sr := NewSearch()
		err := Decode(test.query, sr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "RtrvAttr", sr.RtrvAttr, test.RtrvAttr)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "Facets", sr.Facets, test.Facets)
		if err != nil {
			t.Error(err)
		}
		err = sliceTest(num, "FacetFltr", sr.FacetFltr, test.FacetFltr)
		if err != nil {
			t.Error(err)
		}
		if sr.Index != test.Index {
			t.Errorf("test %v Index: got %#v, expected %#v\n", num, sr.Index, test.Index)
		}
		if sr.ID != test.ID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", num, sr.ID, test.ID)
		}
		if sr.Path != test.Path {
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, sr.Path, test.Path)
		}
	}
}
