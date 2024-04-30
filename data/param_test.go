package data

import (
	"fmt"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
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
	d, err := newData()
	if err != nil {
		t.Error(err)
	}

	params := param.New()
	params.SrchAttr = []string{"title"}
	params.Facets = []string{"tags"}

	var docs []*doc.Doc
	for _, da := range d.data {
		docs = append(docs, doc.New().SetMapping(doc.NewMapping(params)).SetData(da))
	}
	want := 7252
	if len(docs) != want {
		t.Errorf("indexed %v docs, expected %v\n", len(docs), want)
	}
}

func TestSearchFields(t *testing.T) {
	d, err := newData()
	if err != nil {
		t.Error(err)
	}

	params := param.New()
	params.SrchAttr = []string{"title", "comments"}
	params.Facets = []string{"tags"}

	var docs []*doc.Doc
	for _, da := range d.data {
		docs = append(docs, doc.New().SetMapping(doc.NewMapping(params)).SetData(da))
	}

	var ids []int
	for _, doc := range docs {
		if doc.SearchAllFields("falling fish") {
			ids = append(ids, doc.ID)
		}
	}
	ids = lo.Uniq(ids)
	//fmt.Printf("res %#v\n", ids)
	if len(ids) != 2 {
		t.Errorf("got %v results, expected %v\n", len(ids), 2)
	}
}

func TestSearchFacets(t *testing.T) {
	d, err := newData()
	if err != nil {
		t.Error(err)
	}

	params := param.New()
	params.SrchAttr = []string{"title"}
	params.Facets = []string{"tags"}

	var docs []*doc.Doc
	for _, da := range d.data {
		docs = append(docs, doc.New().SetMapping(doc.NewMapping(params)).SetData(da))
	}

	var ids []int
	for _, doc := range docs {
		ids = append(ids, doc.SearchFacets("litrpg")...)
	}

	if len(ids) != 2 {
		t.Errorf("got %v results, expected %v\n", len(ids), 2)
	}
}

func newData() (*Data, error) {
	d := NewData()
	d.AddFile(`../testdata/ndbooks.ndjson`)

	err := d.decodeNDJSON()

	return d, err
}
