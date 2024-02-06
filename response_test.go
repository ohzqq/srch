package srch

import (
	"testing"
)

func TestResponseParams(t *testing.T) {
	idx := newTestIdx()
	res := idx.Search("page=0&query")
	m := res.StringMap()
	hits, ok := m[Hits].([]map[string]any)
	if !ok {
		t.Error("wrong")
	}
	hpp := res.HitsPerPage()
	if len(hits) != hpp {
		t.Errorf("got %d, expected %d\n", len(hits), hpp)
	}
	title, ok := hits[0][DefaultField].(string)
	if !ok {
		t.Errorf("not a string")
	}
	tw := "Cross & Crown"
	if title != tw {
		t.Errorf("sorting err, got %s, expected %s\n", title, tw)
	}

	res = idx.Search("page=1&query")
	m = res.StringMap()
	hits, ok = m[Hits].([]map[string]any)
	if !ok {
		t.Error("wrong")
	}
	hpp = res.HitsPerPage()
	if len(hits) != hpp {
		t.Errorf("got %d, expected %d\n", len(hits), hpp)
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
