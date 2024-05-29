package endpoint

import (
	"slices"
	"testing"
)

func TestRoutes(t *testing.T) {
	want := wantEndpoints
	got := Routes
	if !slices.Equal(got, want) {
		t.Errorf("got %v endoints, wanted %v\n", got, want)
	}
}

var tests = []string{
	"/indexes",
	"/indexes/audiobooks",
	"/indexes/audiobooks/browse",
	"/indexes/audiobooks/tags",
	"/indexes/audiobooks/query",
	"/indexes/audiobooks/settings",
	"/indexes/audiobooks/facets",
	"/indexes/audiobooks/facets/tags",
	"/indexes/audiobooks/facets/tags/query",
}

func TestWildcards(t *testing.T) {
	var wildcards = []string{"audiobooks", "tags"}
	for i, end := range Endpoints {
		want := tests[i]
		got := end.SetWildcards(wildcards...)
		if got != want {
			t.Errorf("endpoint %v\ngot %v path, want %v\n", end.Route(), got, want)
		}
	}
}

var wantEndpoints = []string{
	"/indexes",
	"/indexes/{indexName}",
	"/indexes/{indexName}/browse",
	"/indexes/{indexName}/{objectID}",
	"/indexes/{indexName}/query",
	"/indexes/{indexName}/settings",
	"/indexes/{indexName}/facets",
	"/indexes/{indexName}/facets/{facetName}",
	"/indexes/{indexName}/facets/{facetName}/query",
}
