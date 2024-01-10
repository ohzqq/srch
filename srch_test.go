package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/audible"
)

func testVals() url.Values {
	vals := make(url.Values)
	vals.Add("tags", "abo")
	vals.Add("tags", "dnr")
	//vals.Add("authors", "Alice Winters")
	//vals.Add("authors", "Amy Lane")
	//vals.Add("q", "fish")
	return vals
}

func TestParseValues(t *testing.T) {
	vals, err := ParseValues(testQueryString)
	if err != nil {
		t.Error(err)
	}
	if len(vals["tags"]) != 2 {
		t.Errorf("got %d, expected 2", len(vals["tags"]))
	}
}

func TestFileSrc(t *testing.T) {
	src := FileSrc(testData)
	testFilterQueryString(t, src)
}

func TestSliceSrc(t *testing.T) {
	src := SliceMapSrc(books)
	testFilterQueryString(t, src)
}

func TestReadDataSrc(t *testing.T) {
	f, err := os.Open(testData)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	src := ReaderSrc(f)
	testFilterQueryString(t, src)
}

func TestFuzzySearch(t *testing.T) {
	data := make([]map[string]any, len(books))
	for i, book := range books {
		data[i] = map[string]any{"title": book["title"]}
	}
	m := FuzzyFind(data, "fish")
	fmt.Printf("gen text %v\n", m.Len())
}

func TestGenericFullTextSearch(t *testing.T) {
	data := make([]string, len(books))
	for i, book := range books {
		data[i] = book["title"].(string)
	}
	ft := New()
	res := ft.Search("fish", StringSliceSrc(data))
	if res.Len() != 8 {
		t.Errorf("got %d, expected %d\n", res.Len(), 8)
	}
}

func testFilterQueryString(t *testing.T, src DataSrc) {
	res := idx.Index(src()).Filter(testQueryString)
	if len(res.Data) != 2 {
		t.Errorf("got %d, expected %d\n", len(res.Data), 2)
	}
}

func TestFilterData(t *testing.T) {
	d := Filter(books, idx.FacetFields(), testVals())
	if len(d) != 384 {
		t.Errorf("got %d, expected %d\n", len(d), 384)
	}
}

func TestFullTextSearch(t *testing.T) {
	res := idx.Search("fish", SliceMapSrc(books))
	if len(res.Data) != 8 {
		t.Errorf("got %d, expected 8\n", len(res.Data))
	}
}

func TestFullTextSearchChoose(t *testing.T) {
	t.SkipNow()
	res := idx.Search("fish", SliceMapSrc(books))
	if len(res.Data) != 8 {
		t.Errorf("got %d, expected 8\n", len(res.Data))
	}
	sel, err := res.Choose()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sel.Len())
}

func TestSearchAndFilter(t *testing.T) {
	res := idx.Search("fish", SliceMapSrc(books))
	if len(res.Data) != 8 {
		t.Errorf("got %d, expected 8\n", len(res.Data))
	}

	q := "authors=Amy+Lane"
	f := res.Filter(q)
	//fmt.Printf("facets %+v\n", idx.Facets()[0])
	fmt.Println(len(f.Data))
}

func TestAudibleSearch(t *testing.T) {
	a := New(
		WithSearch(audibleSrch),
		WithTextFields([]string{"Title"}),
	)
	println("audible search")

	res := a.Search("amy lane fish")
	fmt.Printf("num res %d\n", len(res.Data))
	fmt.Printf("num res %v\n", res.searchableFields)

	//for i := 0; i < res.Len(); i++ {
	//  println(res.String(i))
	//}

	//res.Print()
}

func audibleSrc(q string) DataSrc {
	return func() []map[string]any {
		return audibleApi(q)
	}
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
	//fmt.Printf("products %v\n", r.Products)
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
