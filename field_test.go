package srch

import (
	"fmt"
	"log"
	"testing"
)

var fieldSortParams = []string{}

func TestFieldSort(t *testing.T) {
	var sortErr error
	alpha := libCfgStr + "&sortFacetValuesBy=alpha"
	idx, err := New(alpha)
	if err != nil {
		log.Fatal(err)
	}
	tags := idx.GetFacet("tags")
	tags.Order = "asc"
	sorted := tags.SortTokens()
	switch tags.Order {
	case "desc":
		if sorted[0].Label != "zombies" {
			sortErr = fmt.Errorf("alpha: %s (%d)\n", sorted[0].Label, sorted[0].Count())
		}
	case "asc":
		if sorted[0].Label != "abo" {
			sortErr = fmt.Errorf("alpha: %s (%d)\n", sorted[0].Label, sorted[0].Count())
		}
	}

	count := libCfgStr + "&sortFacetValuesBy=count"
	idx, err = New(count)
	if err != nil {
		log.Fatal(err)
	}
	tags = idx.GetFacet("tags")
	tags.Order = "asc"
	sorted = tags.SortTokens()
	switch tags.Order {
	case "desc":
		if sorted[0].Label != "dnr" {
			sortErr = fmt.Errorf("count: %s (%d)\n", sorted[0].Label, sorted[0].Count())
		}
	case "asc":
		if sorted[0].Label != "courting" {
			sortErr = fmt.Errorf("count: %s (%d)\n", sorted[0].Label, sorted[0].Count())
		}
	}
	if sortErr != nil {
		t.Error(sortErr)
	}
}
