package idx

import (
	"testing"
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
	for _, test := range paramTests {
		idx, err := Open(test.query)
		if err != nil {
			t.Fatal(err)
		}
		err = checkIdxName(idx, "index")
		if err != nil {
			t.Fatal(err)
		}
		err = checkSrchAttr(idx)
	}
}
