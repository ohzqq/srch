package srch

import (
	"net/http"
	"net/http/httptest"

	"github.com/ohzqq/srch/endpoint"
	"github.com/spf13/cast"
)

func NewSrv() *http.Server {
	mux := http.NewServeMux()
	for _, end := range endpoint.Endpoints {
		switch end {
		case endpoint.Root:
			// list indexes
			mux.HandleFunc(end.Get(), http.HandlerFunc(Indexes))
		case endpoint.Idx:
			// search index
			mux.HandleFunc(end.Get(), http.HandlerFunc(IdxReq))
			// add object
			mux.HandleFunc(end.Post(), http.HandlerFunc(IdxObject))
		case endpoint.IdxBrowse:
			// return all objects
			mux.HandleFunc(end.Get(), http.HandlerFunc(IdxBrowse))
			// return all objects
			mux.HandleFunc(end.Post(), http.HandlerFunc(IdxBrowse))
		case endpoint.IdxObject:
			// get object by id
			mux.HandleFunc(end.Get(), http.HandlerFunc(IdxObject))
			// add or replace object by id
			mux.HandleFunc(end.Put(), http.HandlerFunc(IdxObject))
			// delete object by id
			mux.HandleFunc(end.Del(), http.HandlerFunc(IdxObject))
		case endpoint.IdxQuery:
			// search index via form
			mux.HandleFunc(end.Post(), http.HandlerFunc(IdxQuery))
		case endpoint.IdxSettings:
			// get settings
			mux.HandleFunc(end.Get(), http.HandlerFunc(IdxSettings))
			// set settings
			mux.HandleFunc(end.Put(), http.HandlerFunc(IdxSettings))
		case endpoint.Facets:
			// list facets for index
			mux.HandleFunc(end.Get(), http.HandlerFunc(Facets))
		case endpoint.Facet:
			// return all facet values
			mux.HandleFunc(end.Get(), http.HandlerFunc(Facet))
		case endpoint.FacetQuery:
			// Search facet values
			mux.HandleFunc(end.Get(), http.HandlerFunc(FacetQuery))
			mux.HandleFunc(end.Post(), http.HandlerFunc(FacetQuery))
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

func getWildCards(r *http.Request) []string {
	cards := []string{
		r.PathValue(endpoint.IdxName),
	}

	name := r.PathValue(endpoint.FacetName)
	if name != "" {
		cards = append(cards, name)
		return cards
	}

	id := r.PathValue(endpoint.ObjectID)
	if id == "" {
		cards = append(cards, id)
		return cards
	}

	return cards
}

func getIdxName(r *http.Request) (string, bool) {
	name := r.PathValue(endpoint.IdxName)
	if name == "" {
		return name, false
	}
	return name, true
}

func getFacetName(r *http.Request) (string, bool) {
	name := r.PathValue(endpoint.FacetName)
	if name == "" {
		return name, false
	}
	return name, true
}

func getObjectID(r *http.Request) (int, bool) {
	id := r.PathValue(endpoint.ObjectID)
	if id == "" {
		return 0, false
	}
	return cast.ToInt(id), true
}
