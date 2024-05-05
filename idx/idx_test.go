package idx

import (
	"testing"

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
