package srch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"testing"
)

func TestServer(t *testing.T) {
	//mux := Mux()
	ts := OfflineSrv()
	defer ts.Close()

	ts.URL += "/test/poot"
	res, err := http.PostForm(ts.URL, url.Values{})
	if err != nil {
		t.Error(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%s\n", greeting)
}

func TestRoutes(t *testing.T) {
	want := wantEndpoints
	got := Routes
	if !slices.Equal(got, want) {
		t.Errorf("got %v endoints, wanted %v\n", got, want)
	}
}

var wantEndpoints = []string{
	"/indexes",
	"/indexes/{indexName}",
	"/indexes/{indexName}/browse",
	"/indexes/{indexName}/query",
	"/indexes/{indexName}/settings",
	"/indexes/{indexName}/facets",
	"/indexes/{indexName}/facets/{facetName}",
	"/indexes/{indexName}/facets/{facetName}/query",
}
