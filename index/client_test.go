package index

import (
	"testing"
)

const hareTestPath = `/home/mxb/code/srch/testdata/hare`
const hareTestURL = `file://home/mxb/code/srch/testdata/hare`
const hareTestQuery = `?url=file://home/mxb/code/srch/testdata/hare/`

func TestHareDisk(t *testing.T) {
	_, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%#v\n", idx.Database)
}

func TestHareDiskTbls(t *testing.T) {
	client, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}
	names := client.TableNames()
	for _, n := range names {
		println(n)
	}
}

func TestDefaultIndex(t *testing.T) {
	_, err := New("")
	if err != nil {
		t.Error(err)
	}
}

func TestGetIdx(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	idx, err := c.GetIdx(defaultTbl)
	if err != nil {
		t.Error(err)
	}
	if idx.Index != defaultTbl {
		t.Errorf("got %v, wanted %v\n", idx.Index, defaultTbl)
	}
}
