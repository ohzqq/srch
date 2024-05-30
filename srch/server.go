package srch

import (
	"net/http"
	"net/http/httptest"

	"github.com/ohzqq/srch/endpoint"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
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
			mux.HandleFunc(end.Get(), http.HandlerFunc(IdxSrch))
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

func ParseHTTPRequest(r *http.Request) *Request {
	r.ParseForm()
	cards := getWildCards(r)

	params := lo.Assign(
		map[string][]string(r.Form),
		map[string][]string(r.PostForm),
		map[string][]string(r.URL.Query()),
	)

	return &Request{
		vals:      params,
		URL:       r.URL,
		method:    r.Method,
		wildCards: cards,
	}
}

func Indexes(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	client, err := req.Client()
	if err != nil {
	}
	idx := maps.Keys(client.Indexes())
	testRes(w, req, idx)
}

func IdxSrch(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
}

func IdxBrowse(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
}

func IdxObject(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
}

func IdxQuery(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
}

func IdxSettings(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
}

func Facets(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
}

func Facet(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
}

func FacetQuery(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	testRes(w, req, nil)
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
