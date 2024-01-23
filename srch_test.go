package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os/exec"
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

func TestParseValues(t *testing.T) {
	t.SkipNow()
	vals, err := ParseValues(testQueryString)
	if err != nil {
		t.Error(err)
	}
	if len(vals["tags"]) != 2 {
		t.Errorf("got %d, expected 2", len(vals["tags"]))
	}
}

func testSearchQueryStrings() map[string]int {
	queries := map[string]int{
		"": 7174,
	}
	v := make(url.Values)

	v.Set(ParamQuery, "heart")
	queries[v.Encode()] = 303

	v.Set(ParamQuery, "")
	queries[v.Encode()] = 7174

	return queries
}

func TestFuzzySearch(t *testing.T) {
	//t.SkipNow()
	test := settingsTestVals[7]
	idx := New(test.query)
	totalBooksErr(idx.Len(), test.query)
	if len(idx.TextFields()) != len(test.want.SearchableAttributes) {
		t.Errorf("%s: got %+v, wanted %+v\n", test.query, len(idx.TextFields()), len(test.want.SearchableAttributes))
	}

	for q, want := range testSearchQueryStrings() {
		m := idx.Search(q)
		if m.NbHits != want {
			t.Errorf("%s: num res %d, expected %d \n", q, m.Len(), want)
		}
	}
}

func TestFullTextSearch(t *testing.T) {
	t.SkipNow()
	idx = NewIndex(testValuesCfg, WithFullText())
	//idx.Index(books)
	res := idx.SearchIndex("fish")
	if len(res.Data) != 8 {
		t.Errorf("got %d, expected 8\n", len(res.Data))
	}
	//for _, facet := range idx.Facets() {
	//  for _, item := range facet.Items {
	//    fmt.Printf("%s: %d\n", item.Value, item.Count)
	//  }
	//}

}

func parseValueTest(t *testing.T, q string) {
	_, err := ParseValues(q)
	if err != nil {
		t.Error(err)
	}
}

func TestFilterQueryString(t *testing.T) {
	t.SkipNow()
	q := "series=#gaymers"
	parseValueTest(t, q)

	idx.Index(books)
	res := idx.Filter(q)
	if len(res.Data) != 2 {
		t.Errorf("got %d, expected %d\n", len(res.Data), 2)
	}
	d, err := json.Marshal(res)
	if err != nil {
		t.Error(err)
	}
	println(string(d))
}

func TestFilterData(t *testing.T) {
	t.SkipNow()
	idx.Index(books)
	d := FilterData(books, idx.Facets(), testVals())
	if len(d) != 384 {
		t.Errorf("got %d, expected %d\n", len(d), 384)
	}
}

func TestAudibleSearch(t *testing.T) {
	t.SkipNow()

	q := "field=Title"
	a := NewIndex(
		q,
	)
	res := a.SearchIndex("amy lane fish")

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

type testQ string

func (q testQ) String() string {
	return string(q)
}

//func TestCDB(t *testing.T) {
//  t.SkipNow()
//  s := NewSearch(testS)
//  //err := s.Get()
//  sel, err := s.Get("litrpg")
//  if err != nil {
//    t.Error(err)
//  }
//  fmt.Printf("%#v\n", sel)
//}

//func TestTUI(t *testing.T) {
//  t.SkipNow()
//  s := NewSearch(testS)
//  //err := s.Get()
//  sel, err := s.Get("litrpg")
//  if err != nil {
//    t.Error(err)
//  }
//  fmt.Printf("%#v\n", sel)
//}

func cdbSearch(t *testing.T) []byte {
	//cdb := exec.Command("echo", `angst`)

	cdb := exec.Command("calibredb", testS.cmd...)
	//println(cdb.String())

	out, err := cdb.Output()
	if err != nil {
		t.Error(err)
	}
	return out
}

type testResult map[string]any

func (s testSearcher) Search(queries string) ([]any, error) {
	testS.cmd = append(testS.cmd, "-s", queries)
	cdb := exec.Command("calibredb", testS.cmd...)
	println(cdb.String())

	out, err := cdb.Output()
	if err != nil {
		return nil, err
	}

	var res []testResult
	err = json.Unmarshal(out, &res)
	if err != nil {
		return nil, err
	}

	var items []any
	for _, r := range res {
		items = append(items, r)
	}

	return items, nil
}

func (r testResult) String() string {
	return r["title"].(string)
}
