package idx

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

const (
	facetParamStr   = `facets=tags,authors,series,narrators`
	facetParamSlice = `facets=tags&facets=authors&facets=series&facets=narrators`
	srchAttrParam   = "searchableAttributes=title"
	queryParam      = `query=fish`
	sortParam       = `sortBy=title&order=desc`
	filterParam     = `facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`
	uidParam        = `uid=id`
)

func TestNewDB(t *testing.T) {
	idx := New()
	err := checkIdxName(idx, "index")
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
		err = checkIdxName(idx, "index")
		if err != nil {
			t.Error(err)
		}
		err = checkAttrs(param.SrchAttr.String(), idx.Params.SrchAttr, test.want.SrchAttr)
		if err != nil {
			t.Errorf("\nparams: %v\ntest num %v: %v\n", test.query, i, err)
		}
	}
}
