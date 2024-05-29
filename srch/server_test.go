package srch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/ohzqq/srch/endpoint"
)

func TestServer(t *testing.T) {
	//mux := Mux()
	ts := OfflineSrv()
	defer ts.Close()

	println(endpoint.IdxName)

	ts.URL += "/indexes/audiobooks"
	res, err := http.PostForm(ts.URL, url.Values{})
	if err != nil {
		t.Error(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%s\n", greeting)
}
