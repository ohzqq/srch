package srch

import (
	"fmt"
	"net/url"
	"testing"
)

func TestResponseFacets(t *testing.T) {
	idx := newTestIdx()

	jq := `{"facetFilters":["authors:amy lane", ["tags:romance"]],"facets":["authors","narrators","series","tags"],"maxValuesPerFacet":200,"page":0,"query":""}`
	params := ParseSearchParamsJSON(jq)
	res := idx.Search(params)

	if res.HasFilters() {
		filters := res.Params.Get(FacetFilters)
		println(filters)
		res = res.Filter(filters)
	}

	for _, facet := range res.facets {
		fmt.Printf("%+v\n", facet.Len())
	}
}

func TestParseFilterJSON(t *testing.T) {
	tf := `["authors:amy lane",["tags:romance"]]`
	jq := `{"facetFilters":["authors:amy lane", ["tags:romance"]],"facets":["authors","narrators","series","tags"],"maxValuesPerFacet":200,"page":0,"query":""}`
	parsed := parseSearchParamsJSON(jq)

	parsedFilter := parsed.Get(FacetFilters)
	if tf != parsedFilter {
		t.Errorf("parsed %s, og %s\n", parsedFilter, tf)
	}

	esc, err := url.QueryUnescape(parsedFilter)
	if err != nil {
		t.Error(err)
	}
	if tf != esc {
		t.Errorf("parsed %s, og %s\n", esc, tf)
	}
}

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
