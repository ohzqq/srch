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

func TestIdxInsertData(t *testing.T) {
	test := func(idx *Idx) error {
		rc, err := idx.openData()
		if err != nil {
			t.Error(err)
		}
		defer rc.Close()

		err = idx.Batch(rc)
		if err != nil {
			t.Error(err)
		}

		ct := idx.DataContentType()
		switch ct {
		case NdJSON:
			//f, err := os.Open()
			//println("need to idx ndjson to mem table")
		case JSON:
			println("need to idx to mem table")
		case Hare:
			println("need to load hare table")
		}
		return nil
	}
	req, err := NewRequest(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series&uid=id&name=audiobooks`)
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
	//runIdxTests(t, test)
}

func TestDataContentType(t *testing.T) {
	test := func(idx *Idx) error {
		ct := idx.DataContentType()
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
	runIdxTests(t, test)
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
