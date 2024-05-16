package param

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

type searchTest struct {
	pt
	*Search
}

var srchTests = []searchTest{
	searchTest{
		pt: pt{
			query: ``,
		},
		Search: &Search{
			RtrvAttr: []string{"*"},
			Paramz: &Paramz{
				Index: "default",
			},
		},
	},
	searchTest{
		pt: pt{
			query: `?path=file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
		},
		Search: &Search{
			RtrvAttr: []string{"*"},
			Paramz: &Paramz{
				Index: "default",
				Path:  `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
			},
		},
	},
	searchTest{
		pt: pt{
			query: `?path=file://home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish`,
		},
		Search: &Search{

			RtrvAttr: []string{"*"},
			Query:    "fish",
			Paramz: &Paramz{
				Index: "default",
				Path:  `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
			},
		},
	},
	searchTest{
		pt: pt{
			query: `?path=file://home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors`,
		},
		Search: &Search{
			Query:    "fish",
			RtrvAttr: []string{"title", "tags", "authors"},
			Paramz: &Paramz{
				Index: "default",
				Path:  `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
			},
		},
	},
	searchTest{
		pt: pt{
			query: `?path=file://home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title&attributesToRetrieve=tags&attributesToRetrieve=authors`,
		},
		Search: &Search{
			Query:    "fish",
			RtrvAttr: []string{"title", "tags", "authors"},
			Paramz: &Paramz{
				Index: "default",
				Path:  `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
			},
		},
	},
	searchTest{
		pt: pt{
			query: `?path=file://home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators`,
		},
		Search: &Search{
			Query:    "fish",
			RtrvAttr: []string{"title", "tags", "authors"},
			Facets:   []string{"tags", "authors", "series", "narrators"},
			Paramz: &Paramz{
				Path:  `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
				Index: "default",
			},
		},
	},
	searchTest{
		pt: pt{
			query: `?path=file://home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&facetFilters=["authors:amy lane"]`,
		},
		Search: &Search{
			Query:     "fish",
			RtrvAttr:  []string{"title", "tags", "authors"},
			Facets:    []string{"tags", "authors", "series", "narrators"},
			FacetFltr: []string{"[\"authors:amy lane\"]"},
			Paramz: &Paramz{
				Index: "default",
				Path:  `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
			},
		},
	},
	searchTest{
		pt: pt{
			query: `?path=file://home/mxb/code/srch/testdata/ndbooks.ndjson&query=fish&attributesToRetrieve=title,tags,authors&facets=tags,authors,series,narrators&page=3&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]&index=audiobooks`,
		},
		Search: &Search{
			RtrvAttr:  []string{"title", "tags", "authors"},
			Page:      3,
			Query:     "fish",
			SortBy:    "title",
			Order:     "desc",
			Facets:    []string{"tags", "authors", "series", "narrators"},
			FacetFltr: []string{"[\"authors:amy lane\", [\"tags:romance\", \"tags:-dnr\"]]"},
			Paramz: &Paramz{
				Path:  `file://home/mxb/code/srch/testdata/ndbooks.ndjson`,
				Index: "audiobooks",
			},
		},
	},
}

func TestDecodeSearchStr(t *testing.T) {
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
		if num == 6 {
			w := []any{"authors:amy lane"}
			f := sr.FacetFilters()
			for i, ff := range f {
				if w[i] != ff {
					t.Errorf("got %v filters, expected %v\n", ff, w[i])
				}
			}
		}
		if num == 7 {
			w := []any{
				"authors:amy lane",
				[]any{"tags:romance", "tags:-dnr"},
			}
			f := sr.FacetFilters()
			if len(f) != 2 {
				t.Errorf("got %v filters, expected %v\n", len(f), 2)
			}
			if f[0] != w[0] {
				t.Errorf("got %v filters, expected %v\n", f[0], w[0])
			}
			err = sliceTest(num, "filters", cast.ToStringSlice(f[1]), cast.ToStringSlice(w[1]))
			if err != nil {
				t.Error(err)
			}
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

func TestDecodeSearchURL(t *testing.T) {
	for num, test := range srchTests {
		sr := NewSearch()
		err := Decode(test.url(), sr)
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
		if num == 6 {
			w := []any{"authors:amy lane"}
			f := sr.FacetFilters()
			for i, ff := range f {
				if w[i] != ff {
					t.Errorf("got %v filters, expected %v\n", ff, w[i])
				}
			}
		}
		if num == 7 {
			w := []any{
				"authors:amy lane",
				[]any{"tags:romance", "tags:-dnr"},
			}
			f := sr.FacetFilters()
			if len(f) != 2 {
				t.Errorf("got %v filters, expected %v\n", len(f), 2)
			}
			if f[0] != w[0] {
				t.Errorf("got %v filters, expected %v\n", f[0], w[0])
			}
			err = sliceTest(num, "filters", cast.ToStringSlice(f[1]), cast.ToStringSlice(w[1]))
			if err != nil {
				t.Error(err)
			}
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
			fmt.Printf("%#v\n", test.url().Query().Get("path"))
			t.Errorf("test %v Path: got %#v, expected %#v\n", num, sr.Path, test.Path)
		}
	}
}
