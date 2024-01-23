package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/RoaringBitmap/roaring"
)

func TestUnmarshalQueryParams(t *testing.T) {
	params := &Query{}
	err := json.Unmarshal(testParamsBytes(), params)
	if err != nil {
		t.Error(err)
	}
	filters, err := params.GetFacetFilters()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", filters)
	filtersTests(filters, t)
}

func TestParseFilterString(t *testing.T) {
	enc := testComboFilterEnc()
	filters, err := DecodeFilter(enc)
	if err != nil {
		t.Error(err)
	}
	filtersTests(filters, t)
}

func filtersTests(filters *Filters, t *testing.T) {
	if len(filters.Con) != 1 {
		t.Errorf("got %d conjunctive filters, expected %d\n", len(filters.Con), 1)
	}
	if len(filters.Dis) != 2 {
		t.Errorf("got %d disjunctive filters, expected %d\n", len(filters.Dis), 2)
	}
}

func TestMarshalFilter(t *testing.T) {
	combo := testComboFilter()
	var c []any
	err := json.Unmarshal([]byte(combo), &c)
	if err != nil {
		t.Error(err)
	}
}

func TestFiltering(t *testing.T) {
	test := settingsTestVals[7]
	idx := New(test.query)
	totalBooksErr(idx.Len(), test.query)

	for q, want := range testSearchQueryStrings() {
		req := NewQuery(q)
		if req.HasFilters() {
			ids, err := Filter(idx.Bitmap(), idx.Fields, q)
			if err != nil {
				t.Error(err)
			}
			if len(ids) != want {
				t.Errorf("query %s:\ngot %d, expected %d\n", q, len(ids), want)
			}
		}
	}
}

func TestSearchAndFilter(t *testing.T) {
	test := settingsTestVals[7]
	idx := New(test.query)

	vals := make(url.Values)
	vals.Set(ParamQuery, "heart")

	//total := 7174

	vals.Set(ParamFacetFilters, `["authors:amy lane"]`)

	afterFilter := 5
	//afterFilter := 58

	result := idx.Search(vals.Encode())
	if n := result.NbHits(); n != afterFilter {
		t.Errorf("got %d, expected %d\n", n, afterFilter)
	}

	d, err := json.Marshal(result)
	if err != nil {
		t.Error(err)
	}
	println(string(d))

	idx.res = roaring.New()
	f := idx.Filter(`["authors:amy lane"]`)
	if n := f.res.GetCardinality(); n != 58 {
		t.Errorf("got %d, expected %d\n", n, 58)
	}
	//vals.Set(ParamFacetFilters, `["authors:amy lane", [ "tags:romance"], "tags:-dnr"]`)
}

func testSearchFilterStrings() map[string]int {
	queries := map[string]int{}
	v := make(url.Values)

	v.Set(ParamFacetFilters, `["authors:amy lane"]`)
	queries[v.Encode()] = 58

	v.Set(ParamFacetFilters, `["authors:amy lane", ["tags:romance"]]`)
	queries[v.Encode()] = 26

	v.Set(ParamFacetFilters, `["authors:amy lane", ["tags:romance"], "tags:-dnr"]`)
	queries[v.Encode()] = 22

	return queries
}

var plainFilters = []string{
	`"authors:amy lane", ["series:fish"]`,
	`["tag:grumpy/sunshine","tag:-enemies to lovers"]`,
}

var encodedFilters = []string{
	`%22authors%3Aamy+lane%22`,
	`%5B%22tag%3Agrumpy%2Fsunshine%22%2C+%22tag%3Aenemies+to+lovers%22%5D`,
}

func testOrFilter() string {
	return fmt.Sprint("[", plainFilters[1], "]")
}

func testEncOrFilter() string {
	return url.QueryEscape(testOrFilter())
}

func testComboFilter() string {
	f := fmt.Sprintf("[%s,%s]", plainFilters[0], plainFilters[1])
	println(f)
	return f
}

func testComboFilterEnc() string {
	return url.QueryEscape(testComboFilter())
}

func testAndFilter() string {
	return fmt.Sprint("[", plainFilters[0], "]")
}

func testEncAndFilter() string {
	return url.QueryEscape(testAndFilter())
}

func testParamsBytes() []byte {
	d, err := json.Marshal(requestParams())
	if err != nil {
		log.Fatal(err)
	}
	return d
}
