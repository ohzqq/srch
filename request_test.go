package srch

import (
	"fmt"
	"testing"

	"github.com/ohzqq/srch/param"
)

var testFields = []string{
	param.SrchAttr,
	param.FullText,
	param.DataFile,
	param.DataDir,
	param.Facets,
	param.Page,
	param.Query,
	param.SortBy,
	param.Order,
	param.FacetFilters,
}

func TestNewRequest(t *testing.T) {
	for i := 0; i < 3; i++ {
		req := NewRequest().
			FullText(`testdata/poot.bleve`).
			UID("id").
			Query("fish").
			Page(i)
			//HitsPerPage(5)

		res, err := idx.Search(req.String())
		if err != nil {
			t.Fatal(err)
		}

		err = searchErr(res.NbHits(), 37, res.Params.Query)
		if err != nil {
			t.Error(err)
		}

		hits := res.Hits()
		//fmt.Printf("%#v\n", res.nbHits[0]["title"])
		if len(hits) > 0 {
			title := hits[0]["title"].(string)
			switch i {
			case 0:
				want := "Fish on a Bicycle"
				if title != want {
					fmt.Printf("got %s, wanted %s\n", title, want)
				}
			case 1:
				want := "Hide and Seek"
				if title != want {
					fmt.Printf("got %s, wanted %s\n", title, want)
				}
			}
		}
	}
}
