package srch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

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

func TestOrFilter(t *testing.T) {
	t.SkipNow()
	println(testOrFilter())
	println(testEncOrFilter())
	println(testAndFilter())
	println(testEncAndFilter())
	println(testComboFilter())
	enc := testComboFilterEnc()
	println(enc)
	dec, err := url.QueryUnescape(enc)
	if err != nil {
		t.Error(err)
	}
	println(dec)
}

func TestParseFilterString(t *testing.T) {
	enc := testComboFilterEnc()
	filters, err := ParseFilterJSONString(enc)
	if err != nil {
		t.Error(err)
	}
	if len(filters.And) != 1 {
		t.Errorf("got %d conjunctive filters, expected %d\n", len(filters.And), 1)
	}
	if len(filters.Or) != 2 {
		fmt.Printf("%+v\n", filters)

		t.Errorf("got %d disjunctive filters, expected %d\n", len(filters.Or), 2)
	}
}

func TestFilterToString(t *testing.T) {
	filters := testFilterStruct.String()
	if filters != testComboFilter() {
		t.Errorf("got %v, expected %s\n", filters, testComboFilter())
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

func testOrFilter() string {
	return fmt.Sprint("[", plainFilters[1], "]")
	//return plainFilters[1]
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
	//return plainFilters[0]
}

func testEncAndFilter() string {
	return url.QueryEscape(testAndFilter())
}
