package srch

import (
	"net/http"
	"net/http/cgi"
	"testing"
)

func cgiHandler(w http.ResponseWriter, r *http.Request) {
	handler := cgi.Handler{Path: "/mnt/c/Users/nina/code/cgi/cgi_test.sh"}
	handler.ServeHTTP(w, r)
}

func TestCgi(t *testing.T) {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}
