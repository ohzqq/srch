package srch

import (
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/audible"
)

const testQueryString = `tags=grumpy/sunshine&tags=enemies+to+lovers`
const testSearchString = `q=amy+lane`

func testVals() url.Values {
	vals := make(url.Values)
	vals.Add("tags", "abo")
	vals.Add("tags", "dnr")
	//vals.Add("authors", "Alice Winters")
	//vals.Add("authors", "Amy Lane")
	//vals.Add(QueryField, "fish")
	return vals
}

func testSearchQueryStrings() map[string]int {
	queries := map[string]int{
		"": 7174,
	}
	v := make(url.Values)

	v.Set(Query, "heart")
	queries[v.Encode()] = 303

	v.Set(Query, "")
	queries[v.Encode()] = 7174

	return queries
}

func TestFuzzySearch(t *testing.T) {
	idx := newTestIdx()

	for q, want := range testSearchQueryStrings() {
		m := idx.Search(q)
		if m.NbHits() != want {
			t.Errorf("%s: num res %d, expected %d \n", q, m.NbHits(), want)
		}
	}
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

	vals := make(url.Values)
	vals.Set(Query, "fish")

	res := idx.Search(vals.Encode())
	if h := res.NbHits(); h != 8 {
		t.Errorf("get %d, expected %d\n", h, 8)
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

type testSearcher struct {
	cmd []string
}

var testS = testSearcher{
	cmd: []string{"list", "--with-library", "http://localhost:8888/#audiobooks", "--username", "churkey", "--password", "<f:/home/mxb/.dotfiles/misc/calibre.txt>", "--limit", "2", "--for-machine"},
}
