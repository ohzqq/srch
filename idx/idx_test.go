package idx

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

func TestDefaultIdx(t *testing.T) {
	idx := NewIdx()

	i := 0

	err := checkIdxName(i, idx.Params.IndexName)
	if err != nil {
		t.Fatal(err)
	}

}

func TestOpenIdx(t *testing.T) {
	for i, test := range paramTests {
		idx, err := Open(test.query)
		if err != nil {
			t.Fatalf("test no %v: %v\n", i, err)
		}

		err = checkIdxPath(i, idx.Params.Path)
		if err != nil {
			t.Fatal(err)
		}

		if idx.Params.Path != "" {
			ids, err := idx.DB.IDs("index")
			if err != nil {
				t.Error(err)
			}

			if len(ids) != 7251 {
				t.Errorf("test no %v: got %v, wanted %v\n", i, len(ids), 7251)
			}
		}
	}
}

func TestConfigureIdx(t *testing.T) {
	for i, test := range paramTests {
		idx := Init(test.query)

		err := checkIdxName(i, idx.Params.IndexName)
		if err != nil {
			t.Error(err)
		}

		err = checkAttrs(i, param.SrchAttr, idx.Params.SrchAttr)
		if err != nil {
			t.Errorf("\nparams: %v\n%v\n", test.query, err)
		}

		err = checkAttrs(i, param.Facets, idx.Params.Facets)
		if err != nil {
			t.Errorf("\nparams: %v\n%v\n", test.query, err)
		}

		err = checkAttrs(i, param.FacetAttr, idx.Params.FacetAttr)
		if err != nil {
			t.Errorf("\nparams: %v\n%v\n", test.query, err)
		}
	}
}

//func TestSearchIdx(t *testing.T) {
//  tests := []map[string][]string{
//    map[string][]string{
//      "dragon": []string{"title"},
//    },
//    map[string][]string{
//      "omega": []string{"title"},
//    },
//    map[string][]string{
//      "dragon omega": []string{"title"},
//    },
//    map[string][]string{
//      "dragon": []string{"comments"},
//    },
//    map[string][]string{
//      "omega": []string{"comments"},
//    },
//    map[string][]string{
//      "dragon omega": []string{"comments"},
//    },
//    map[string][]string{
//      "dragon": []string{"title", "comments"},
//    },
//    map[string][]string{
//      "omega": []string{"title", "comments"},
//    },
//    map[string][]string{
//      "dragon omega": []string{"title", "comments"},
//    },
//  }

//  want := []int{
//    74,
//    97,
//    1,
//    185,
//    328,
//    23,
//    200,
//    345,
//    23,
//  }
//  dsk, err := db.NewDiskStore(testHareDskDir)
//  if err != nil {
//    t.Fatal(err)
//  }

//  for i, test := range tests {
//    for kw, attrs := range test {
//      params := param.New()
//      params.SrchAttr = attrs

//      m := NewMappingFromParams(params)
//      db, err := db.New()
//      if err != nil {
//        t.Error(err)
//      }
//      err = db.Init(dsk)
//      if err != nil {
//        t.Fatal(err)
//      }

//      ids, err := db.Search(kw)
//      if err != nil {
//        t.Error(err)
//      }
//      if res := len(ids); res != want[i] {
//        t.Errorf("kw %s, attrs %v: got %v results, wanted %v\n", kw, attrs, res, want[i])
//      }
//    }
//  }
//}
