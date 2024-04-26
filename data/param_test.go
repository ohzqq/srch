package data

import (
	"fmt"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
)

var testQuerySettings = []string{
	"",
	"searchableAttributes=",
	"searchableAttributes=title&fullText=../testdata/poot.bleve",
	"searchableAttributes=title&dataDir=../testdata/data-dir",
	"attributesForFaceting=tags,authors,series,narrators",
	"attributesForFaceting=tags,authors,series,narrators&dataFile=../testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators",
	"searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators",
	`searchableAttributes=title&dataFile=../testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
}

const hareTestDB = `../testdata/hare`

type Record struct {
	ID     int
	Fields map[string]*bloom.BloomFilter
}

func TestNewClient(t *testing.T) {
	t.SkipNow()
	tp := `../testdata/data-dir`

	data := New("dir", tp)
	d, err := data.Decode()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("data %#v\n", len(d))
}

func TestNewData(t *testing.T) {
	d := NewData()
	d.AddFile(`../testdata/ndbooks.ndjson`)

	err := d.decodeNDJSON()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("data %#v\n", len(d.data))
	fmt.Printf("ids %#v\n", len(d.ids))
}
