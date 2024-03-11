package srch

import (
	"fmt"
	"log"
	"net/url"
	"slices"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/audible"
	"github.com/ohzqq/srch/blv"
)

func TestNewBlvSearch(t *testing.T) {
	var i Indexer
	i = blv.Open(`testdata/poot.bleve`)
	bits, err := i.Search("fish")
	if err != nil {
		t.Error(err)
	}
	if h := bits.GetCardinality(); h != 8 {
		t.Errorf("got %d hits, expected %d\n", h, 8)
	}

	idx := newTestIdxCfg("&fullText=testdata/poot.bleve")
	res := idx.Search("query=fish")
	if h := res.NbHits(); h != 8 {
		t.Errorf("got %d hits, expected %d\n", h, 8)
	}
}

func testSearchQueryStrings() map[string]int {
	queries := map[string]int{
		"": numBooks,
	}
	v := make(url.Values)

	v.Set(Query, "fish")
	queries[v.Encode()] = 303

	v.Set(Query, "")
	queries[v.Encode()] = numBooks

	return queries
}

func TestFuzzySearch(t *testing.T) {
	idx := newTestIdx()

	err := srchTest(idx, 56)
	if err != nil {
		t.Error(err)
	}

}

func srchTest(idx *Idx, want int) error {
	err := searchErr(idx, numBooks, "")
	if err != nil {
		return err
	}

	err = searchErr(idx, want, "query=fish")
	if err != nil {
		return err
	}
	return nil
}

func searchErr(idx *Idx, want int, q string) error {
	res := idx.Search(q)
	err := intErr(res.NbHits(), want, q)
	if err != nil {
		return err
	}
	op, err := url.QueryUnescape(idx.Params.String())
	if err != nil {
		return err
	}
	rp, err := url.QueryUnescape(res.Params.String())
	if err != nil {
		return err
	}
	if op != rp {
		for key, val := range res.Params.Search {
			if key == "facets" {
				return nil
			}
			has := idx.Params.Has(key)
			if !has {
				return fmt.Errorf("doesn't have key %v\n", key)
			}
			o := idx.Params.Search[key]
			if !slices.Equal(val, o) {
				return fmt.Errorf("idx params: %s\nres params: %s\n", op, rp)
			}
		}
	}
	return nil
}

//func TestFullTextSearch(t *testing.T) {
//  idx := newTestIdx()
//  println(idx.Len())

//  res := idx.Search("query=fish")

//  println(res.NbHits())

//  cfg := libCfgStr + "&indexPath=" + blevePath

//  println(cfg)
//  idx = newTestIdxCfg("&indexPath=" + blevePath)
//  res = idx.Search("query=fish")

//  println(res.NbHits())
//  //if ana := idx.GetAnalyzer(); ana != TextAnalyzer {
//  //t.Errorf("get %s, expected %s\n", ana, TextAnalyzer)
//  //}

//  //err = srchTest(idx, 8)
//  //if err != nil {
//  //t.Error(err)
//  //}
//}

func parseValueTest(t *testing.T, q string) {
	_, err := ParseValues(q)
	if err != nil {
		t.Error(err)
	}
}

func TestAudibleSearch(t *testing.T) {
	t.SkipNow()

	q := "field=Title"
	a, err := New(q)
	if err != nil {
		t.Error(err)
	}
	res := a.Search("amy lane fish")

	println("audible search")

	//res := a.Search("amy lane fish")
	fmt.Printf("num res %d\n", res.Len())

}

func audibleSrch(q string) []map[string]any {
	return audibleApi(q)
}

func audibleApi(q string) []map[string]any {
	s := audible.NewSearch(q)
	r, err := audible.Products().Search(s)
	if err != nil {
		log.Fatal(err)
	}
	var sl []map[string]any
	for _, p := range r.Products {
		a := make(map[string]any)
		mapstructure.Decode(p, &a)
		sl = append(sl, a)
	}
	fmt.Printf("products %v\n", r.Products)
	return sl
}

var testCalibreStr = []string{"list", "--with-library", "http://localhost:8888/#audiobooks", "--username", "churkey", "--password", "<f:/home/mxb/.dotfiles/misc/calibre.txt>", "--limit", "2", "--for-machine"}
