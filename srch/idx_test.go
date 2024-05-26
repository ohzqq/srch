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

func TestInitIdx(t *testing.T) {
	runTests(t, testInitIdx)
}

func testInitIdx(_ int, req reqTest) error {
	client, err := req.Client()
	if err != nil {
		return err
	}

	idx, err := client.FindIdx(client.IndexName())
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", idx.db)

	return nil
}
