package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"testing"
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
	if len(filters.Dis) != 1 {
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
