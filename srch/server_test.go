package srch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

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

func testHandler(w http.ResponseWriter, r *http.Request) {
	req := NewReq(r)
	fmt.Fprintf(w, "%#v", req)
}
