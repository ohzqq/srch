package srch

import (
	"fmt"
	"slices"
	"testing"

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

func TestFuzzySearch(t *testing.T) {
	req := NewRequest().
		SetRoute(testDataFile).
		SrchAttr("title").
		Query("fish")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Log(err)
	}

	if res.nbHits() != 56 {
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
		SetRoute(testDataDir).
		UID("id").
		Facets("tags", "authors", "narrators", "series").
		AndFilter("authors:amy lane").
		SrchAttr("title")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Fatal(err)
	}

	if 58 != res.NbHits {
		t.Errorf("got hits %d, expected hits %#v\n", 58, res.NbHits)
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
		SetRoute(testDataDir).
		UID("id").
		SrchAttr("title").
		Facets("tags", "authors", "narrators", "series").
		Query("fish")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Fatal(err)
	}

	for _, facet := range res.FacetFields {
		for _, tok := range facet.Keywords() {
			ids := lo.ToAnySlice(tok.Items())
			rel := FilterDataByID(res.results, ids, res.Params.UID)
			i := 0
			for _, r := range rel {
				if facet.Attribute != "tags" {
					continue
				}
				f, ok := r[facet.Attribute]
				if ok {
					vals := cast.ToStringSlice(f)
					if slices.Contains(vals, tok.Label) != true {
						t.Errorf("hit %v does not contain val %s", f, tok.Label)
					}
				}
				i++
			}
			//if i != len(rel) {
			//  t.Errorf("got %d hits with val, expected %d\n", i, len(rel))
			//}
		}
	}
}

func TestNewRequest(t *testing.T) {
	for i := 0; i < 3; i++ {
		req := NewRequest().
			SetRoute(testBlvPath).
			UID("id").
			Query("fish").
			Facets("tags").
			Page(i)
			//HitsPerPage(5)

		res, err := idx.Search(req.String())
		if err != nil {
			t.Fatal(err)
		}

		err = searchErr(res.NbHits, 37, res.Params.Query)
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

	req := NewRequest()

	req.SetRoute(testDataFile)

	for _, key := range param.SettingParams {
		switch key {
		case param.SrchAttr:
			val := viper.GetStringSlice(key)
			req.SrchAttr(val...)
		case param.FacetAttr:
			val := viper.GetStringSlice(key)
			req.SrchAttr(val...)
		case param.SortAttr:
			val := viper.GetStringSlice(key)
			req.SrchAttr(val...)
		case param.UID:
			val := viper.GetString(key)
			req.UID(val)
		case param.Format:
			val := viper.GetString(key)
			req.Format(val)
		}
	}
	for _, key := range param.SearchParams {
		switch key {
		case param.SortFacetsBy:
			val := viper.GetString(key)
			req.SortFacetsBy(val)
		case param.Facets:
			val := viper.GetStringSlice(key)
			req.Facets(val...)
		case param.Filters:
			val := viper.GetString(key)
			req.Filters(val)
		case "or":
			val := viper.GetStringSlice("or")
			req.OrFilter(val...)
		case "and":
			val := viper.GetStringSlice("and")
			req.AndFilter(val...)
		case param.RtrvAttr:
			val := viper.GetStringSlice(key)
			req.RtrvAttr(val...)
		case param.Page:
			val := viper.GetInt(key)
			req.Page(val)
		case param.HitsPerPage:
			val := viper.GetInt(key)
			req.HitsPerPage(val)
		case param.Query:
			val := viper.GetString(key)
			req.Query(val)
		case param.SortBy:
			val := viper.GetString(key)
			req.SortBy(val)
		case param.Order:
			val := viper.GetString(key)
			req.Order(val)
		}
	}
	fmt.Printf("flags %s\n", req)
}

func searchErr(got int, want int, q string) error {
	if got != want {
		return fmt.Errorf("query %s got %d results, wanted %d\n", q, got, want)
	}
	return nil
}
