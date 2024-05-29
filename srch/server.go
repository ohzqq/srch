package srch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

const (
	segmentIndexes   = `/indexes`
	segmentIndexName = `{indexName}`
	segmentQuery     = `query`
	segmentBrowse    = `browse`
	segmentSettings  = `settings`
	segmentFacets    = `facets`
	segmentFacetName = `{facetName}`
)

const (
	apiBase     = "/indexes"
	IdxEndpoint = "/indexes/{indexName}"
)

var (
	IdxSrchEndpoint   = filepath.Join(IdxEndpoint, "query")
	IdxBrowseEndpoint = filepath.Join(IdxEndpoint, "browse")
	IdxCfgEndpoint    = filepath.Join(IdxEndpoint, "settings")
	FacetsEndpoint    = filepath.Join(IdxEndpoint, "facets")
	FacetEndpoint     = filepath.Join(FacetsEndpoint, "{facetName}")
	FacetSrchEndpoint = filepath.Join(FacetEndpoint, "query")
)

type endpoint struct {
	Idx         string
	IdxBrowse   string
	IdxQuery    string
	IdxSettings string
	Facets      string
	Facet       string
	FacetQuery  string
}

var Endpoint = endpoint{
	Idx:         segmentIndexes,
	IdxBrowse:   filepath.Join(segmentIndexes, segmentIndexName, segmentBrowse),
	IdxQuery:    filepath.Join(segmentIndexes, segmentIndexName, segmentQuery),
	IdxSettings: filepath.Join(segmentIndexes, segmentIndexName, segmentSettings),
	Facets:      filepath.Join(segmentIndexes, segmentIndexName, segmentFacets),
	Facet:       filepath.Join(segmentIndexes, segmentIndexName, segmentFacets, segmentFacetName),
	FacetQuery:  filepath.Join(segmentIndexes, segmentIndexName, segmentFacets, segmentFacetName, segmentQuery),
}

func NewSrv() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	mux.HandleFunc("POST /test/{path}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "path %s\n", r.PathValue("path"))
	}))
	srv := &http.Server{
		Handler: mux,
	}
	return srv
}

func OfflineSrv() *httptest.Server {
	ts := httptest.NewUnstartedServer(nil)
	ts.Config = NewSrv()
	ts.Start()
	return ts
}

var Endpoints = []string{
	apiBase,
	IdxEndpoint,
	IdxSrchEndpoint,
	IdxBrowseEndpoint,
	IdxCfgEndpoint,
	FacetsEndpoint,
	FacetEndpoint,
	FacetSrchEndpoint,
}
