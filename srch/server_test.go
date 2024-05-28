package srch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestServer(t *testing.T) {
	//mux := Mux()
	ts := OfflineSrv()
	defer ts.Close()

	ts.URL += "/test/poot"
	res, err := http.PostForm(ts.URL, url.Values{})
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
