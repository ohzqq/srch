package srch

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ohzqq/srch/endpoint"
)

func NewSrv() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint.Root.Get(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	mux.HandleFunc(endpoint.Idx.Post(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "route %v\npath %s\n", endpoint.Idx.Post(), r.PathValue(endpoint.IdxName))
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
