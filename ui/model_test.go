package ui

import (
	"fmt"
	"testing"

	"github.com/ohzqq/srch"
)

const testData = `../testdata/data-dir/audiobooks.json`
const testCfgFile = `../testdata/config.json`

func testSearch(t *testing.T) *srch.Index {
	idx := newIdx()
	res := idx.Search("fish")
	if len(res.Data) != 8 {
		t.Fatalf("got %d, expected 8\n", len(res.Data))
	}
	return res
}

func TestChoose(t *testing.T) {
	t.SkipNow()
	res := testSearch(t)
	sel, err := Choose(res)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sel.Len())
}

func TestRefineFacet(t *testing.T) {
	t.SkipNow()
	res := testSearch(t)
	auth, err := res.GetFacet("authors")
	if err != nil {
		t.Error(err)
	}
	sel := FilterFacet(auth)
	fmt.Printf("res items %v\n", sel)
	filtered := res.Filter(sel)
	fmt.Printf("res filtered %v\n", filtered.Len())
	//println(filtered.Len())
}

func TestFacets(t *testing.T) {
	idx := newIdx()
	auth, err := idx.GetFacet("tags")
	if err != nil {
		t.Error(err)
	}
	sel := FilterFacet(auth)
	fmt.Printf("res items %v\n", sel)
	filtered := idx.Filter(sel)
	fmt.Printf("res filtered %v\n", filtered.Len())
}

func newIdx() *srch.Index {
	idx := srch.OldNew(srch.WithCfg(testCfgFile))
	return idx.Index(srch.FileSrc(testData))
}
