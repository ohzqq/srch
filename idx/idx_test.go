package idx

import "testing"

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
	if idx.Name != "index" {
		t.Error("idx.Name isn't index")
	}
}

func TestOpenIdx(t *testing.T) {
	idx, err := Open(srchAttrParam)
	if err != nil {
		t.Fatal(err)
	}
	if idx.Name != "index" {
		t.Error("idx.Name isn't index")
	}
}
