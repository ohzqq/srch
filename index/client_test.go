package index

import (
	"slices"
	"testing"
)

func TestHareDisk(t *testing.T) {
	_, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHareDiskTbls(t *testing.T) {
	client, err := New(hareTestQuery)
	if err != nil {
		t.Fatal(err)
	}
	got := client.TableNames()
	want := defTbls
	if !slices.Equal(got, want) {
		t.Errorf("got %v tables, wanted %v\n", got, want)
	}
}

func TestDefaultClient(t *testing.T) {
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
