package srch

import (
	"fmt"
	"testing"
)

var bleveSearchTests = []string{
	`blv/testdata/poot.bleve?searchableAttributes=title&uid=id`,
	`blv/testdata/poot.bleve?searchableAttributes=title&uid=id&facets=tags,authors,narrators,series`,
	`blv/testdata/poot.bleve?searchableAttributes=*&uid=id&facets=tags,authors,narrators,series`,
}

var blvfacetCount = map[string]int{
	"tags":      218,
	"authors":   1612,
	"series":    1740,
	"narrators": 1428,
}

var fishfacetCount = map[string]int{
	"tags":      39,
	"authors":   29,
	"series":    22,
	"narrators": 26,
}

func TestFuzzySearch(t *testing.T) {
	req := NewRequest().
		SetRoute("file/testdata/nddata/ndbooks.ndjson").
		SrchAttr("title").
		Query("fish")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Log(err)
	}

	if res.NbHits() != 56 {
		fmt.Printf("route %s, path %s\n", req.Route, req.Path)
		t.Fatal()
	}
}

func TestBleveSearchAll(t *testing.T) {
	for i := 0; i < len(bleveSearchTests); i++ {
		q := bleveSearchTests[i]
		query := ""
		if i == 1 {
			query = "&query=fish"
		}
		res, err := idx.Search(q + query)
		if err != nil {
			t.Error(err)
		}
		got := len(res.hits)
		want := 7252
		if query == "&query=fish" {
			want = 37
			//want = len(res)
		}
		err = searchErr(got, want, query)
		if err != nil {
			t.Error(err)
		}

		if res.Facets != nil {
			for _, facet := range res.Facets {
				if query == "&query=fish" {
					if num, ok := fishfacetCount[facet.Attribute]; ok {
						if num != facet.Len() {
							t.Errorf("%v got %d, expected %d \n", facet.Attribute, facet.Len(), num)
						}
					} else {
						t.Errorf("attr %s not found\n", facet.Attribute)
					}
				} else {
					if num, ok := blvfacetCount[facet.Attribute]; ok {
						if num != facet.Len() {
							t.Errorf("%v got %d, expected %d \n", facet.Attribute, facet.Len(), num)
						}
					} else {
						t.Errorf("attr %s not found\n", facet.Attribute)
					}
				}
			}
		}
	}
}

func TestBleveFacets(t *testing.T) {
	t.SkipNow()
	for i := 2; i < len(bleveSearchTests); i++ {
		q := bleveSearchTests[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		if !idx.isMem {
			t.Errorf("not bleve")
		}
		query := ""
		if i == 1 {
			query = "fish"
		}
		res, err := idx.Search(query)
		if err != nil {
			t.Error(err)
		}

		for _, facet := range res.Facets {
			if num, ok := blvfacetCount[facet.Attribute]; ok {
				if num != facet.Len() {
					t.Errorf("%v got %d, expected %d \n", facet.Attribute, facet.Len(), num)
				}
			} else {
				t.Errorf("attr %s not found\n", facet.Attribute)
			}
		}
	}
}

func searchErr(got int, want int, q string) error {
	if got != want {
		return fmt.Errorf("query %s got %d results, wanted %d\n", q, got, want)
	}
	return nil
}
