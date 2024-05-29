package srch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
