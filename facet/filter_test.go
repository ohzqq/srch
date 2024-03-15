package facet

import (
	"net/url"
	"testing"

	"github.com/ohzqq/srch/param"
)

var filterStrs = []filterStr{
	filterStr{
		want:  2240,
		query: `data=../testdata/nddata/ndbooks.ndjson&facets=tags&facetFilters=["tags:dnr"]`,
	},
	filterStr{
		want:  384,
		query: `data=../testdata/nddata/ndbooks.ndjson&facets=tags&facetFilters=["tags:dnr", "tags:abo"]`,
	},
	filterStr{
		want:  33,
		query: `data=../testdata/nddata/ndbooks.ndjson&facets=tags&facetFilters=["tags:-dnr", "tags:abo"]`,
	},
	filterStr{
		want:  33,
		query: `data=../testdata/nddata/ndbooks.ndjson&facets=tags&facetFilters=["tags:abo", "tags:-dnr"]`,
	},
	filterStr{
		want:  2273,
		query: `data=../testdata/nddata/ndbooks.ndjson&facets=tags&facetFilters=[["tags:dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  5397,
		query: `data=../testdata/nddata/ndbooks.ndjson&facets=tags&facetFilters=[["tags:-dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  5397,
		query: `data=../testdata/nddata/ndbooks.ndjson&facets=tags&facetFilters=[[ "tags:abo", "tags:-dnr"]]`,
	},
}

type filterStr struct {
	query string
	want  int
}

type filterVal struct {
	vals url.Values
	want int
}

func TestFilterStrings(t *testing.T) {
	for _, f := range filterStrs {
		p, err := param.Parse(f.query)
		if err != nil {
			t.Error(err)
		}

		data, err := loadData()
		if err != nil {
			t.Error(err)
		}

		facets, err := New(data, p.FacetSettings)
		if err != nil {
			t.Error(err)
		}

		if num := facets.Len(); num != f.want {
			t.Errorf("query %s:\ngot %d results, wanted %d\n", f.query, num, f.want)
		}
	}

}

func testSearchFilterStrings() []filterVal {
	//queries := make(map[int]url.Values)
	var queries []filterVal

	queries = append(queries, filterVal{
		want: 58,
		vals: url.Values{
			"data":                  []string{"../testdata/nddata/ndbooks.ndjson"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["authors:amy lane"]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 26,
		vals: url.Values{
			"data":                  []string{"../testdata/nddata/ndbooks.ndjson"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["authors:amy lane", ["tags:romance"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 41,
		vals: url.Values{
			//"uid": []string{"id"},
			"data":                  []string{"../testdata/nddata/ndbooks.ndjson"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 0,
		vals: url.Values{
			"data":                  []string{"../testdata/nddata/ndbooks.ndjson"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["tags:abo", "tags:dnr", "tags:horror"]`,
			},
		},
	})

	return queries
}
