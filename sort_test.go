package srch

import (
	"fmt"
	"testing"

	"github.com/ohzqq/srch/param"
	"github.com/spf13/viper"
)

func TestSortAttr(t *testing.T) {
	var sortAttrs = []string{
		"title",
		"title:desc",
		"title:asc:string",
		"added:asc:int",
	}

	for _, test := range sortAttrs {
		s := NewSort(test)
		fmt.Printf("sort %#v\n", s)
	}
}

func TestSortRequest(t *testing.T) {
	req := NewRequest().
		SetRoute(param.Dir.String()).
		UID("id").
		Facets("tags", "authors", "narrators", "series").
		SortAttr("added:asc:int").
		SetPath(testDataDir)

	res, err := idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	println(viper.GetInt(param.HitsPerPage.Snake()))
	fmt.Printf("%#v\n", res.Hits[0])
}
