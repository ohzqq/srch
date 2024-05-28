package srch

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	mux := Mux()
	ts := httptest.NewUnstartedServer(mux)
	ts.Config = NewSrv()
	ts.Start()
	defer ts.Close()

	ts.URL += "/test/poot"
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%s", greeting)
}
