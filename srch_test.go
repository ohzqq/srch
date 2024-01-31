package srch

import (
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/audible"
)

func testSearchQueryStrings() map[string]int {
	queries := map[string]int{
		"": 7174,
	}
	v := make(url.Values)

	v.Set(Query, "fish")
	queries[v.Encode()] = 303

	v.Set(Query, "")
	queries[v.Encode()] = 7174

	return queries
}

func TestFuzzySearch(t *testing.T) {
	idx := newTestIdx()

	err := srchTest(idx, 56)
	if err != nil {
		t.Error(err)
	}

}

func srchTest(idx *Index, want int) error {
	err := searchErr(idx, 7174, "")
	if err != nil {
		return err
	}

	err = searchErr(idx, want, "query=fish")
	if err != nil {
		return err
	}
	return nil
}

func searchErr(idx *Index, want int, q string) error {
	res := idx.Search(q)
	err := intErr(res.NbHits(), want, q)
	if err != nil {
		return err
	}
	op := idx.Params.String()
	rp := res.Params.String()
	if op != rp {
		return fmt.Errorf("idx params %s\nres params%s\n", op, rp)
	}
	return nil
}

func TestFuzzyFieldSearch(t *testing.T) {
	idx := newTestIdx()

	facet := idx.GetFacet("authors")
	if total := facet.Len(); total != len(facet.GetLabels()) {
		t.Errorf("got %d, expected %d\n", len(facet.GetLabels()), facet.Len())
	}
}

func TestFullTextSearch(t *testing.T) {
	cfg := libCfgStr + "&fullText"
	idx, err := New(cfg)
	if err != nil {
		t.Error(err)
	}

	if ana := idx.GetAnalyzer(); ana != TextAnalyzer {
		t.Errorf("get %s, expected %s\n", ana, TextAnalyzer)
	}

	err = srchTest(idx, 8)
	if err != nil {
		t.Error(err)
	}
}

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
