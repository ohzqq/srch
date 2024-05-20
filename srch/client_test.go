package srch

import (
	"testing"
)

type clientTest struct {
	*Client
}

func TestClientMem(t *testing.T) {
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Fatal(err)
		}
		test, err := req.clientTest()
		if err != nil {
			t.Fatal(err)
		}
		want := req.clientWant(i)
		println(test.IndexName())
		println(want.IndexName())

	}
}

func TestClientDisk(t *testing.T) {
}

func TestClientNet(t *testing.T) {
}
