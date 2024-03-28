package srch

import (
	"fmt"
	"slices"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var bleveSearchTests = []string{
	blvRoute(uidParam),
	blvRoute(srchAttrParam, uidParam, facetParamSlice),
	blvRoute("searchableAttributes=*", uidParam, facetParamStr),
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

func TestDecode(t *testing.T) {
	p := map[string]interface{}{"facets": []interface{}{"tags", "authors", "narrators", "series"}, "hits_per_page": 25, "max_values_per_facet": 25, "order": "desc", "searchable_attributes": []interface{}{"title"}, "sort_by": "added", "uid": "id"}

	params := param.New()
	err := mapstructure.Decode(p, params)
	if err != nil {
		t.Error(err)
	}
	println(params.String())
}

func TestFilterAuth(t *testing.T) {

	params := "?searchableAttributes=title&facets=tags,authors,series,narrators&hitsPerPage=25&order=desc&searchableAttributes=title&sortBy=added&uid=id"
	idx, err := New(params)
	if err != nil {
		t.Error(err)
	}
	books := loadData(t)
	idx.Batch(books)
	err = totalBooksErr(7252, params)
	if err != nil {
		t.Error(err)
	}

	q := `?facetFilters=%5B%22authors%3Aandrew+grey%22%5D&facets=authors&facets=tags&facets=narrators&facets=series&hitsPerPage=25&order=desc&searchableAttributes=title&sortBy=added&uid=id`

	res, err := idx.Search(q)
	if err != nil {
		t.Log(err)
	}

	err = searchErr(res.NbHits, 99, q)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("first hits %#v\n", res.Hits[0])
}

func TestFuzzySearch(t *testing.T) {
	req := NewRequest().
		SetRoute(param.File.String()).
		SetPath(testDataFile).
		SrchAttr("title").
		Query("fish")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Log(err)
	}
	println(req.String())

	err = searchErr(res.nbHits(), 56, req.String())
	if err != nil {
		t.Error(err)
	}
}

func TestBleveSearchAll(t *testing.T) {
	for i := 0; i < len(bleveSearchTests); i++ {
		q := bleveSearchTests[i]
		query := ""
		if i == 1 {
			query = "&query=fish"
		}
		sq := q + query
		res, err := idx.Search(sq)
		if err != nil {
			t.Error(err)
		}
		got := len(res.results)
		want := 7252
		if query == "&query=fish" {
			want = 37
			//want = len(res)
		}
		err = searchErr(got, want, query)
		if err != nil {
			t.Error(err)
		}

		if res.FacetFields != nil {
			for _, facet := range res.FacetFields {
				if query == "&query=fish" {
					if num, ok := fishfacetCount[facet.Attribute]; ok {
						if num != facet.Len() {
							t.Errorf("q %s\n%v got %d, expected %d \n", sq, facet.Attribute, facet.Len(), num)
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
	//t.SkipNow()
	for i := 2; i < len(bleveSearchTests); i++ {
		q := bleveSearchTests[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		query := ""
		if i == 1 {
			query = "fish"
		}
		res, err := idx.Search(query)
		if err != nil {
			t.Error(err)
		}

		for _, facet := range res.FacetFields {
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

var facetCount = map[string]int{
	"tags":      33,
	"authors":   5,
	"series":    13,
	"narrators": 19,
}

func TestFacetFilters(t *testing.T) {
	req := NewRequest().
		SetRoute(param.Dir.String()).
		SetPath(testDataDir).
		UID("id").
		Facets("tags", "authors", "narrators", "series").
		AndFilter("authors:amy lane").
		SrchAttr("title")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	if 58 != res.NbHits {
		t.Errorf("got hits %d, expected hits %#v\n", res.NbHits, 58)
	}

	if res.Facets != nil {
		for _, facet := range res.FacetFields {
			if num, ok := facetCount[facet.Attribute]; ok {
				if num != facet.Len() {
					t.Errorf("%v got %d, expected %d \n", facet.Attribute, facet.Len(), num)
				}
			} else {
				t.Errorf("attr %s not found\n", facet.Attribute)
			}
		}
	}

	name := "Amy Lane"
	for _, hit := range res.Hits {
		if m, ok := hit["authors"]; ok {
			auth := cast.ToStringSlice(m)
			if !slices.Contains(auth, name) {
				t.Errorf("got authors %v, should include %s\n", auth, name)
			}
		}
	}

}

func TestFacets(t *testing.T) {
	req := NewRequest().
		SetRoute(param.Dir.String()).
		SetPath(testDataDir).
		UID("id").
		SrchAttr("title").
		Facets("tags", "authors", "narrators", "series").
		Query("fish")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	dnr := res.FilterByFacetValue("tags", "dnr")
	if len(dnr) != 22 {
		t.Errorf("got %d hits with val, expected %d\n", len(dnr), 22)
	}

	facet, err := res.Facets.GetFacet("tags")
	if err != nil {
		t.Error(err)
	}
	for _, tok := range facet.Items {
		ids := lo.ToAnySlice(tok.RelatedTo)
		//fmt.Printf("ids %v\n", ids)
		rel := FilterDataByID(res.results, ids, res.UID)

		if len(rel) != len(ids) {
			t.Errorf("got %d hits with val, expected %d\n", len(rel), len(ids))
		}

		i := 0
		for _, r := range rel {
			f, ok := r[facet.Attribute]
			if ok {
				vals := cast.ToStringSlice(f)
				if slices.Contains(vals, tok.Label) != true {
					t.Errorf("hit %v does not contain val %s", f, tok.Label)
				}
			}
			i++
		}
	}

	//d, err := json.Marshal(facet)
	//if err != nil {
	//t.Error(err)
	//}
	//println(string(d))
}

func TestNewRequest(t *testing.T) {
	for i := 0; i < 3; i++ {
		req := NewRequest().
			SetRoute(param.Blv.String()).
			SetPath(testBlvPath).
			UID("id").
			Query("fish").
			Facets("tags").
			Page(i)
			//HitsPerPage(5)

		res, err := idx.Search(req.String())
		if err != nil {
			t.Error(err)
		}

		err = searchErr(res.NbHits, 37, res.Query)
		if err != nil {
			t.Error(err)
		}

		hits := res.Hits
		//fmt.Printf("%#v\n", res.nbHits[0]["title"])
		if len(hits) > 0 {
			title := hits[0]["title"].(string)
			switch i {
			case 0:
				want := "Fish on a Bicycle"
				if title != want {
					fmt.Printf("got %s, wanted %s\n", title, want)
				}
			case 1:
				want := "Hide and Seek"
				if title != want {
					fmt.Printf("got %s, wanted %s\n", title, want)
				}
			}
		}

		//d, err := json.Marshal(res)
		//if err != nil {
		//t.Error(err)
		//}
		//println(string(d))
	}
}

func TestFlags(t *testing.T) {
	viper.Set(param.Blv.String(), "testdata/poot.bleve")

	req := GetViperParams()

	res, err := idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	err = searchErr(res.NbHits, numBooks, res.Query)
	if err != nil {
		t.Error(err)
	}

}

func searchErr(got int, want int, q string) error {
	if got != want {
		return fmt.Errorf("query %s got %d results, wanted %d\n", q, got, want)
	}
	return nil
}
