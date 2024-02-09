package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"testing"
)

var boolFilterStr = []string{
	`tags:dnr`,
	`tags:dnr AND tags:abo`,
	`tags:dnr OR tags:abo`,
	`NOT tags:dnr AND tags:abo`,
	`NOT tags:dnr OR tags:abo`,
	`tags:dnr AND NOT tags:abo`,
	`tags:dnr OR NOT tags:abo`,
}

func TestMarshalFilter(t *testing.T) {
	combo := testComboFilter()
	var c []any
	err := json.Unmarshal([]byte(combo), &c)
	if err != nil {
		t.Error(err)
	}
}

func TestJSONFilter(t *testing.T) {
	idx := newTestIdx()
	tf := `["authors:amy lane",["tags:romance"]]`

	jq := `{"facetFilters":["authors:amy lane", ["tags:romance"]],"facets":["authors","narrators","series","tags"],"maxValuesPerFacet":200,"page":0,"query":""}`

	parsed := parseSearchParamsJSON(jq)
	//println(parsed.Get(FacetFilters))

	err := searchErr(idx, 806, parsed.Encode())
	if err != nil {
		t.Error(err)
	}

	vals := url.Values{
		FacetFilters: []string{tf},
	}
	err = searchErr(idx, 806, vals.Encode())
	if err != nil {
		t.Error(err)
	}
}

func TestSearchAndFilter(t *testing.T) {
	idx := newTestIdx()

	vals := make(url.Values)
	vals.Set(Query, "heart")

	//total := 7174

	vals.Set(FacetFilters, `["authors:amy lane"]`)

	afterFilter := 5
	//afterFilter := 58

	err := searchErr(idx, afterFilter, vals.Encode())
	if err != nil {
		t.Error(err)
	}
}

func TestNewFilters(t *testing.T) {
	idx := newTestIdx()
	for _, test := range testSearchFilterStrings() {
		params := ParseParams(test.vals)
		filters := params.Filters()
		bits, err := Filter(idx.Bitmap(), idx.facets, filters)
		if err != nil {
			t.Error(err)
		}
		hits := int(bits.GetCardinality())
		if hits != test.want {
			t.Errorf("%#v\ngot %d, expected %d\n", filters, hits, test.want)
		}
	}
}

func TestFilters(t *testing.T) {
	//t.SkipNow()
	idx := newTestIdx()

	for _, test := range testSearchFilterStrings() {
		err := searchErr(idx, test.want, test.vals.Encode())
		if err != nil {
			t.Error(err)
		}
	}
}

type filterStr struct {
	vals url.Values
	want int
}

func testSearchFilterStrings() []filterStr {
	//queries := make(map[int]url.Values)
	var queries []filterStr

	queries = append(queries, filterStr{
		want: 58,
		vals: url.Values{
			FacetFilters: []string{
				`["authors:amy lane"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 806,
		vals: url.Values{
			FacetFilters: []string{
				`["authors:amy lane", ["tags:romance"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 789,
		vals: url.Values{
			FacetFilters: []string{
				`["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 384,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:dnr", "tags:abo"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 1856,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:dnr", "tags:-abo"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 1856,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:-abo", "tags:dnr"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2273,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:dnr", "tags:abo"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2240,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:dnr", "tags:-abo"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2240,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:-abo", "tags:dnr"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 0,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:abo", "tags:dnr", "tags:horror"]`,
			},
		},
	})

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
