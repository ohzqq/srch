package srch

import (
	"slices"
	"testing"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func TestFacets(t *testing.T) {
	req := NewRequest().
		SetRoute(testDataDir).
		UID("id").
		SrchAttr("title").
		Facets("tags", "authors", "narrators", "series").
		Query("fish")

	res, err := idx.Search(req.String())
	if err != nil {
		t.Fatal(err)
	}

	for _, facet := range res.Facets {
		for _, tok := range facet.Keywords() {
			ids := lo.ToAnySlice(tok.Items())
			rel := FilterDataByID(res.results, ids, res.Params.UID)
			i := 0
			for _, r := range rel {
				if facet.Attribute != "tags" {
					continue
				}
				f, ok := r[facet.Attribute]
				if ok {
					vals := cast.ToStringSlice(f)
					if slices.Contains(vals, tok.Label) != true {
						t.Errorf("hit %v does not contain val %s", f, tok.Label)
					}
				}
				i++
			}
			//if i != len(rel) {
			//  t.Errorf("got %d hits with val, expected %d\n", i, len(rel))
			//}
		}
	}
}
