package srch

import (
	"testing"
)

func TestParseFacetSort(t *testing.T) {
	test := []string{
		"tags",
		"tags:count",
		"tags:count:asc",
	}
	for i, str := range test {
		field := &Field{}
		parseAttr(field, str)
		switch i {
		case 0:
			if v := field.Attribute; v != "tags" {
				t.Errorf("wrong attribute %s\n", v)
			}
		case 1:
			if v := field.SortBy; v != "count" {
				t.Errorf("wrong sortby %s\n", v)
			}
		case 2:
			if v := field.Order; v != "asc" {
				t.Errorf("wrong Order %s\n", v)
			}
		}
	}
}

func TestSortFacets(t *testing.T) {
	t.SkipNow()

	q := getNewQuery()
	//query := fmt.Sprintf("%s&%s&%s", testValuesCfg, testQueryString, testSearchString)
	i := NewIndex(q.Encode())
	if i.Len() != 2 {
		println(q.Encode())
		t.Errorf("got %d, expected %d\v", i.Len(), 2)
	}

	tags := i.GetFacet("tags")
	if v := tags.SortBy; v != "count" {
		t.Errorf("wrong sortby %s\n", v)
	}
	if v := tags.Order; v != "desc" {
		t.Errorf("wrong Order %s\n", v)
	}

	authors := i.GetFacet("authors")
	if v := authors.SortBy; v != "label" {
		t.Errorf("wrong sortby %s\n", v)
	}
	if v := authors.Order; v != "asc" {
		t.Errorf("wrong Order %s\n", v)
	}

	//for _, item := range authors.items {
	//  fmt.Printf("%s: count %d\n", item.Label, item.Count())
	//}
}
