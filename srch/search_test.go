package srch

import (
	"fmt"
	"testing"
)

var bleveSearchTests = []string{
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`,
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`,
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id&facets=tags,authors,narrators,series`,
}

var facetCount = map[string]int{
	"tags":      222,
	"authors":   1618,
	"series":    1745,
	"narrators": 1430,
}

func TestBleveSearchAll(t *testing.T) {
	for i := 0; i < len(bleveSearchTests); i++ {
		q := bleveSearchTests[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		if !idx.isBleve {
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
		got := len(res.hits)
		want := 7252
		if query == "fish" {
			want = 8
			//want = len(res)
		}
		err = searchErr(got, want, query)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestBleveFacets(t *testing.T) {
	for i := 2; i < len(bleveSearchTests); i++ {
		q := bleveSearchTests[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		if !idx.isBleve {
			t.Errorf("not bleve")
		}
		query := "fish"
		if i == 1 {
			query = "fish"
		}
		res, err := idx.Search(query)
		if err != nil {
			t.Error(err)
		}
		//got := len(res.hits)
		//want := 7252
		//if query == "fish" {
		//want = 8
		//want = len(res)
		//}
		//err = searchErr(got, want, query)
		//if err != nil {
		//t.Error(err)
		//}

		for _, facet := range res.Facets.Facets {
			if num, ok := facetCount[facet.Attribute]; ok {
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
