package srch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
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
