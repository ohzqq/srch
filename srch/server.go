package srch

import (
	"net/http"
	"net/http/httptest"

	"github.com/ohzqq/srch/endpoint"
)

func NewSrv() *http.Server {
	mux := http.NewServeMux()
	for _, end := range endpoint.Endpoints {
		switch end {
		case endpoint.Root:
			// list indexes
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
		case endpoint.Idx:
			// search index
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
			// add object
			mux.HandleFunc(end.Post(), http.HandlerFunc(testHandler))
		case endpoint.IdxBrowse:
			// return all objects
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
			// return all objects
			mux.HandleFunc(end.Post(), http.HandlerFunc(testHandler))
		case endpoint.IdxObject:
			// get object by id
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
			// add or replace object by id
			mux.HandleFunc(end.Put(), http.HandlerFunc(testHandler))
			// delete object by id
			mux.HandleFunc(end.Del(), http.HandlerFunc(testHandler))
		case endpoint.IdxQuery:
			// search index via form
			mux.HandleFunc(end.Post(), http.HandlerFunc(testHandler))
		case endpoint.IdxSettings:
			// get settings
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
			// set settings
			mux.HandleFunc(end.Put(), http.HandlerFunc(testHandler))
		case endpoint.Facets:
			// list facets for index
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
		case endpoint.Facet:
			// return all facet values
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
		case endpoint.FacetQuery:
			// Search facet values
			mux.HandleFunc(end.Get(), http.HandlerFunc(testHandler))
			mux.HandleFunc(end.Post(), http.HandlerFunc(testHandler))
		}
	}
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
