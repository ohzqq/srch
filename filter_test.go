package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"testing"
)

func TestFilterToString(t *testing.T) {
	filters := testFilterStruct.String()
	if filters != testComboFilter() {
		t.Errorf("got %v, expected %s\n", filters, testComboFilter())
	}
}

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
	if len(filters.And) != 1 {
		t.Errorf("got %d conjunctive filters, expected %d\n", len(filters.And), 1)
	}
	if len(filters.Or) != 2 {
		t.Errorf("got %d disjunctive filters, expected %d\n", len(filters.Or), 2)
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
	`"authors:amy lane"`,
	`["tag:grumpy/sunshine","tag:enemies to lovers"]`,
}

var encodedFilters = []string{
	`%22authors%3Aamy+lane%22`,
	`%5B%22tag%3Agrumpy%2Fsunshine%22%2C+%22tag%3Aenemies+to+lovers%22%5D`,
}

var testFilterStruct = &Filters{
	And: []string{"authors:amy lane"},
	Or:  []string{`tag:grumpy/sunshine`, `tag:enemies to lovers`},
}

func testOrFilter() string {
	return fmt.Sprint("[", plainFilters[1], "]")
}

func testEncOrFilter() string {
	return url.QueryEscape(testOrFilter())
}

func testComboFilter() string {
	return fmt.Sprintf("[%s,%s]", plainFilters[0], plainFilters[1])
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
