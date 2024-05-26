package srch

import (
	"fmt"
	"path/filepath"
	"testing"
)

var dataURLs = []QueryStr{
	QueryStr(`?data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`),
	QueryStr(`?name=audiobooks`),
	QueryStr(`?name=audiobooks&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson`),
}

func TestDataContentType(t *testing.T) {
	runTests(t, testContentType)
}

func testContentType(_ int, req reqTest) error {
	client, err := req.Client()
	if err != nil {
		return err
	}

	idx, err := client.FindIdx(client.IndexName())
	if err != nil {
		return err
	}

	ct := idx.ContentType()
	switch ext := filepath.Ext(idx.dataURL.Path); ext {
	case ".json":
		if ct != JSON {
			return fmt.Errorf("got %v content type, wanted %v\n", ct, JSON)
		}
	case ".ndjson":
		if ct != NdJSON {
			return fmt.Errorf("got %v content type, wanted %v\n", ct, NdJSON)
		}
	case ".hare":
		if ct != Hare {
			return fmt.Errorf("got %v content type, wanted %v\n", ct, Hare)
		}
	}
	return nil
}

func TestIdxTbls(t *testing.T) {
	runIdxTests(t, testHasTbls)
}

func testHasTbls(idx *Idx) error {
	if got := !idx.db.TableExists(idx.idxTblName()); got {
		want := true
		if got != want {
			return fmt.Errorf("got %v for tbl %v, wanted %v\n", got, idx.idxTblName(), want)
		}
	}
	if got := !idx.db.TableExists(idx.dataTblName()); got {
		want := true
		if got != want {
			return fmt.Errorf("got %v for tbl %v, wanted %v\n", got, idx.dataTblName(), want)
		}
	}
	return nil
}

type testIdxFunc func(*Idx) error

func runIdxTests(t *testing.T, test testIdxFunc) {
	for _, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Error(err)
		}

		client, err := req.Client()
		if err != nil {
			t.Error(err)
		}

		idx, err := client.FindIdx(client.IndexName())
		if err != nil {
			t.Error(err)
		}

		err = test(idx)
		if err != nil {
			t.Error(err)
		}
	}
}
