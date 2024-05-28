package srch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/test/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	mux.HandleFunc("/test/{path}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "path %s\n", r.PathValue("path"))
	}))
	return mux
}

func NewSrv() *http.Server {
	mux := Mux()
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
