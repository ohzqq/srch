package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"testing"
)

func TestUnmarshalQueryParams(t *testing.T) {
	params := &Params{}
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

func filtersTests(filters *Filters, t *testing.T) {
	if len(filters.Con) != 2 {
		t.Errorf("got %d conjunctive filters, expected %d\n", len(filters.Con), 2)
	}
	if len(filters.Dis) != 1 {
		t.Errorf("got %d disjunctive filters, expected %d\n", len(filters.Dis), 1)
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

func TestSearchAndFilter(t *testing.T) {
	test := "searchableAttributes=title&attributesForFaceting=tags,authors,series&dataFile=testdata/data-dir/audiobooks.json"
	idx, err := New(test)
	if err != nil {
		t.Error(err)
	}

	vals := make(url.Values)
	vals.Set(Query, "heart")

	//total := 7174

	vals.Set(FacetFilters, `["authors:amy lane"]`)

	afterFilter := 5
	//afterFilter := 58

	result := idx.Search(vals.Encode())
	if n := result.NbHits(); n != afterFilter {
		t.Errorf("got %d, expected %d\n", n, afterFilter)
	}
}

func TestNewFilters(t *testing.T) {
	idx := newTestIdx()
	//fields := idx.Params.newFieldsMap(idx.FacetAttr())
	for _, test := range testSearchFilterStrings() {
		filters := test.vals.Get(FacetFilters)
		bits, err := decodeFilter(idx.Bitmap(), idx.fac, filters)
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
	t.SkipNow()
	test := "searchableAttributes=title&attributesForFaceting=tags,authors,series&dataFile=testdata/data-dir/audiobooks.json"
	idx, err := New(test)
	if err != nil {
		t.Error(err)
	}

	for _, test := range testSearchFilterStrings() {
		res := idx.Search(test.vals.Encode())
		if r := res.NbHits(); r != test.want {
			filters, err := idx.GetFacetFilters()
			if err != nil {
				fmt.Printf("vals %+v\n", test.vals)
				t.Error(err)
			}
			t.Errorf("%#v\ngot %d, expected %d\n", filters, r, test.want)
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
		want: 801,
		vals: url.Values{
			FacetFilters: []string{
				`["authors:amy lane", ["tags:romance"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 784,
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
		want: 1853,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:dnr", "tags:-abo"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 1853,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:-abo", "tags:dnr"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2270,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:dnr", "tags:abo"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2237,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:dnr", "tags:-abo"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2237,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:-abo", "tags:dnr"]]`,
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
