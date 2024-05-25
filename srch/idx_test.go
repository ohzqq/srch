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

	data := NewData(client.DataURL())

	ct := data.ContentType()
	switch ext := filepath.Ext(data.Path); ext {
	case ".json":
		if ct != JSON {
			return fmt.Errorf("got %v content type, wanted %v\n", ct, JSON)
		}
	case ".ndjson":
		if ct != NdJSON {
			return fmt.Errorf("got %v content type, wanted %v\n", ct, NdJSON)
		}
	}
	return nil
}
