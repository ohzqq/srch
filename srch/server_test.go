package srch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/ohzqq/srch/endpoint"
)

func runSrvTests(t *testing.T) {
	for _, query := range TestQueryParams {
		for _, end := range endpoint.Endpoints {
			runSrvTest(t, end, query)
		}
	}
}

func runSrvTest(t *testing.T, end endpoint.Endpoint, query QueryStr) {
	ts := OfflineSrv()
	defer ts.Close()

	v := query.Query()
	name := v.Get("name")
	wPath := end.SetWildcards(name)
	ts.URL += wPath
	ts.URL += query.String()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("%#v\n %#v\n", greeting, err)
	}

	if res.StatusCode != 200 {
		t.Errorf("url %s\ngot %v status, wanted %v\n", res.Request.URL, res.Status, http.StatusOK)
	}

	if res.Request.URL.Path != wPath {
		t.Errorf("got %v path, wanted %v\n", res.Request.URL.Path, wPath)
	}
}

func TestServer(t *testing.T) {
	runSrvTests(t)
}

func TestServerPostForm(t *testing.T) {
	//mux := Mux()
	ts := OfflineSrv()
	defer ts.Close()

	ts.URL += "/indexes/audiobooks?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&primaryKey=id&name=audiobooks"
	form := make(url.Values)
	form.Set("data", "file://home/mxb/code/srch/testdata/data-dir/audiobooks.json")
	res, err := http.PostForm(ts.URL, form)
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

func TestServerIndexes(t *testing.T) {
	//mux := Mux()
	ts := OfflineSrv()
	defer ts.Close()

	ts.URL += "/indexes?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=tags&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series,narrators&primaryKey=id&name=audiobooks"
	res, err := http.Get(ts.URL)
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

func testHandler(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func testRes(w http.ResponseWriter, req *Request, r any) {
	fmt.Fprintf(w, "url %s\nroute %s %s\n%#v", req.URL, req.method, req.Endpoint().Route(), r)
}
