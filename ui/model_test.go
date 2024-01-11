package ui

import (
	"fmt"
	"testing"

	"github.com/ohzqq/srch"
)

var idx = &srch.Index{}

const testData = `testdata/data-dir/audiobooks.json`
const testCfgFile = `testdata/config.json`

func init() {
	idx = srch.New(srch.WithCfg(testCfgFile))
	idx.Index(srch.FileSrc(testData))
}

func TestChooseFacet(t *testing.T) {
	//t.SkipNow()
	res := idx.Search("fish")
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

func TestFullTextSearchChoose(t *testing.T) {
	res := idx.Search("fish")
	if len(res.Data) != 8 {
		t.Errorf("got %d, expected 8\n", len(res.Data))
	}
	sel, err := Choose(res)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sel.Len())
}
