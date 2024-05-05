package idx

import (
	"testing"

	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/param"
)

func TestDefaultIdx(t *testing.T) {
	idx := New()

	i := 0

	err := checkIdxName(i, idx.Name)
	if err != nil {
		t.Fatal(err)
	}

}

func TestOpenIdx(t *testing.T) {
	for i, test := range paramTests {
		idx, err := Open(test.query)
		if err != nil {
			t.Fatal(err)
		}

		err = checkIdxPath(i, idx.Params.Path)
		if err != nil {
			t.Fatal(err)
		}

		if idx.Params.Path != "" {
			println(idx.Params.Path)
		}
	}
}

func TestParams(t *testing.T) {
	for i, test := range paramTests {
		params, err := param.Parse(test.query)
		if err != nil {
			t.Fatal(err)
		}

		err = checkIdxName(i, params.IndexName)
		if err != nil {
			t.Error(err)
		}

		err = checkAttrs(i, param.SrchAttr, params.SrchAttr)
		if err != nil {
			t.Errorf("\nparams: %v\n%v\n", test.query, err)
		}

		err = checkAttrs(i, param.Facets, params.Facets)
		if err != nil {
			t.Errorf("\nparams: %v\n%v\n", test.query, err)
		}

		err = checkAttrs(i, param.FacetAttr, params.FacetAttr)
		if err != nil {
			t.Errorf("\nparams: %v\n%v\n", test.query, err)
		}
	}
}

func TestSearchIdx(t *testing.T) {
	tests := []map[string][]string{
		map[string][]string{
			"dragon": []string{"title"},
		},
		map[string][]string{
			"omega": []string{"title"},
		},
		map[string][]string{
			"dragon omega": []string{"title"},
		},
		map[string][]string{
			"dragon": []string{"comments"},
		},
		map[string][]string{
			"omega": []string{"comments"},
		},
		map[string][]string{
			"dragon omega": []string{"comments"},
		},
		map[string][]string{
			"dragon": []string{"title", "comments"},
		},
		map[string][]string{
			"omega": []string{"title", "comments"},
		},
		map[string][]string{
			"dragon omega": []string{"title", "comments"},
		},
	}

	want := []int{
		74,
		97,
		1,
		185,
		328,
		23,
		200,
		345,
		23,
	}
	dsk, err := db.NewDiskStore(hareTestDB)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		for kw, attrs := range test {
			params := param.New()
			params.SrchAttr = attrs

			m := NewMappingFromParams(params)
			db, err := db.New()
			if err != nil {
				t.Error(err)
			}
			err = db.Init(dsk)
			if err != nil {
				t.Fatal(err)
			}

			ids, err := db.Search(kw)
			if err != nil {
				t.Error(err)
			}
			if res := len(ids); res != want[i] {
				t.Errorf("kw %s, attrs %v: got %v results, wanted %v\n", kw, attrs, res, want[i])
			}
		}
	}
}
