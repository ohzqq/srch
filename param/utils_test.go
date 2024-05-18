package param

import (
	"net/url"
	"strings"
	"testing"
)

const (
	hareTestPath  = `/home/mxb/code/srch/testdata/hare`
	hareTestURL   = `file://home/mxb/code/srch/testdata/hare`
	hareTestQuery = `?url=file://home/mxb/code/srch/testdata/hare/`
)

const (
	dataTestURL = `file://home/mxb/code/srch/testdata/ndbooks.ndjson`
	idxTestFile = `file://home/mxb/code/srch/testdata/hare/audiobooks.json`
)

type queryTest string

type cfgTest struct {
	query string
	*Cfg
}

func (p queryTest) str() string {
	return string(p)
}

func (p queryTest) vals() url.Values {
	v, _ := url.ParseQuery(strings.TrimPrefix(p.str(), "?"))
	return v
}

func (p queryTest) url() *url.URL {
	u, _ := url.Parse(p.str())
	return u
}

func testSrch(t *testing.T, num queryTest, got, want *Search) {
	err := sliceTest(num, "RtrvAttr", got.RtrvAttr, want.RtrvAttr)
	if err != nil {
		t.Error(err)
	}
}

func testCfg(t *testing.T, num queryTest, got, want *Cfg) {
	if got.IndexName() != want.IndexName() {
		t.Errorf("test %v Index: got %#v, expected %#v\n", num, got.IndexName(), want.IndexName())
	}
	if got.Client.UID != want.Client.UID {
		t.Errorf("test %v ID: got %#v, expected %#v\n", num, got.Client.UID, want.Client.UID)
	}
	if got.DataURL().Path != want.DataURL().Path {
		t.Errorf("test %v Path: got %#v, expected %#v\n", num, got.DataURL().Path, want.DataURL().Path)
	}
	if got.DB().Path != want.DB().Path {
		t.Errorf("test %v Path: got %#v, expected %#v\n", num, got.DB().Path, want.DB().Path)
	}
	if got.SrchURL().Path != want.SrchURL().Path {
		t.Errorf("test %v Path: got %#v, expected %#v\n", num, got.SrchURL().Path, want.SrchURL().Path)
	}
}

func testIdx(t *testing.T, num queryTest, got, want *Idx) {
	err := sliceTest(num, "SrchAttr", got.SrchAttr, want.SrchAttr)
	if err != nil {
		t.Error(err)
	}
	err = sliceTest(num, "FacetAttr", got.FacetAttr, want.FacetAttr)
	if err != nil {
		t.Error(err)
	}
	err = sliceTest(num, "SortAttr", got.SortAttr, want.SortAttr)
	if err != nil {
		t.Error(err)
	}
}
