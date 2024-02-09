package srch

import (
	"testing"
)

func TestResponsePagination(t *testing.T) {
	idx := newTestIdx()
	//params := ParseSearchParamsJSON(`{"facets":["authors","narrators","series","tags"],"maxValuesPerFacet":200,"page":0,"query":"","tagFilters":""}`)
	res := idx.Post(`{"facets":["authors","narrators","series","tags"],"maxValuesPerFacet":200,"page":0,"query":"","tagFilters":""}`)
	m := res.StringMap()
	hits, ok := m[Hits].([]map[string]any)
	if !ok {
		t.Error("wrong")
	}

	hpp := res.HitsPerPage()
	if len(hits) != hpp {
		t.Errorf("got %d, expected %d\n", len(hits), hpp)
	}

	if res.Page() != 0 {
		t.Errorf("got %d, expected %d\n", res.Page(), 0)
	}

	title, ok := hits[0][DefaultField].(string)
	if !ok {
		t.Errorf("not a string")
	}
	tw := "Cross & Crown"
	if title != tw {
		t.Errorf("sorting err, got %s, expected %s\n", title, tw)
	}

	params := ParseSearchParamsJSON(`{"facets":["authors","narrators","series","tags"],"maxValuesPerFacet":200,"page":1,"query":"","tagFilters":""}`)
	res = idx.Search(params)
	m = res.StringMap()
	hits, ok = m[Hits].([]map[string]any)
	if !ok {
		t.Error("wrong")
	}
	hpp = res.HitsPerPage()
	if len(hits) != hpp {
		t.Errorf("got %d, expected %d\n", len(hits), hpp)
	}

	if res.Page() != 1 {
		t.Errorf("got %d, expected %d\n", res.Page(), 1)
	}

	title, ok = hits[0][DefaultField].(string)
	if !ok {
		t.Errorf("not a string")
	}
	tw = "DEX"
	if title != tw {
		t.Errorf("sorting err, got %s, expected %s\n", title, tw)
	}
}
