package idx

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

func TestDefaultIdx(t *testing.T) {
	idx := New()
	err := checkIdxName(idx, "index")
	if err != nil {
		t.Fatal(err)
	}

	i := 0

	err = checkAttrs(param.SrchAttr.String(), idx.Params.SrchAttr, paramTests[i].want.SrchAttr)
	if err != nil {
		t.Errorf("\nparams: %v\ntest num %v: %v\n", paramTests[i].query, i, err)
	}

	err = checkAttrs(param.Facets.String(), idx.Params.FacetAttr, paramTests[i].want.FacetAttr)
	if err != nil {
		t.Errorf("\nparams: %v\ntest num %v: %v\n", paramTests[i].query, i, err)
	}
}

func TestOpenIdx(t *testing.T) {
	for i, test := range paramTests {
		idx, err := Open(test.query)
		if err != nil {
			t.Fatal(err)
		}
		err = checkIdxName(idx, "index")
		if err != nil {
			t.Error(err)
		}

		err = checkAttrs(param.SrchAttr.String(), idx.Params.SrchAttr, test.want.SrchAttr)
		if err != nil {
			t.Errorf("\nparams: %v\ntest num %v: %v\n", test.query, i, err)
		}

		err = checkAttrs(param.FacetAttr.String(), idx.Params.FacetAttr, test.want.FacetAttr)
		if err != nil {
			t.Errorf("\nparams: %v\ntest num %v: %v\n", test.query, i, err)
		}
	}
}

func TestParams(t *testing.T) {
	for _, test := range testQuerySettings {
		println(test)
	}
}
