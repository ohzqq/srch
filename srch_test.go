package srch

import (
	"encoding/json"
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
	src := SliceSrc(books)
	testFilterQueryString(t, src)
}

func TestReadDataSrc(t *testing.T) {
	f, err := os.Open(testData)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	src := ReadDataSrc(f)
	testFilterQueryString(t, src)
}

func testFilterQueryString(t *testing.T, src DataSrc) {
	res, err := idx.Filter(src, testQueryString)
	if err != nil {
		t.Error(err)
	}
	if len(res.Data) != 2 {
		t.Errorf("got %d, expected %d\n", len(res.Data), 2)
	}
}

func TestFilterData(t *testing.T) {
	d := Filter(books, idx.Facets(), testVals())
	if len(d) != 384 {
		t.Errorf("got %d, expected %d\n", len(d), 384)
	}
}

func TestSearchResults(t *testing.T) {
	t.SkipNow()
	res := SearchData(
		books,
		idx.Fields,
		FullText(books, "title"),
		"fish",
	)
	data := res.Data
	if len(data) != 8 {
		t.Errorf("got %d, expected 8\n", len(data))
	}
	//for _, f := range res.Facets {
	//  fmt.Printf("attr %s: %+v\n", f.Attribute, f.Items[0])
	//}
}

func TestSearchSrc(t *testing.T) {
	t.SkipNow()
	data := FullTextSearch(books, "fish", "title")
	//data := i.data
	if len(data) != 8 {
		t.Errorf("got %d, expected 8\n", len(data))
	}
}

func TestIdxFilterSearch(t *testing.T) {
	t.SkipNow()
	//vals := testVals()
	//res := idx.Search(vals)

	fn := FuzzySearch(books, "title")
	res := fn("fish")
	i := NewIndex(SliceSrc(res), WithCfg(testCfgFile))
	vals := make(url.Values)
	vals.Set("authors", "amy lane")
	r := i.FilterFacets(vals)
	data := r.GetData()
	if len(data) != 4 {
		t.Errorf("got %d, expected 4", len(data))
	}
}

func TestAudibleSearch(t *testing.T) {
	t.SkipNow()

	a := NewIndex(
		audibleSrc("sporemaggeddon"),
		WithSearch(audibleSrch),
		WithTextFields([]string{"Title"}),
		Interactive,
	)
	res := a.Search("amy lane fish")
	println("audible search")
	println(len(res.Data))

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
