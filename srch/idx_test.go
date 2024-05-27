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

const testIdxReq = QueryStr(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series&uid=id&name=audiobooks`)

const (
	testDocPK = 7312
	testDocID = 7245
)

func TestIdxInsertData(t *testing.T) {
	t.SkipNow()
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
	req, err := NewRequest(`?searchableAttributes=title&db=file://home/mxb/code/srch/testdata/hare&sortableAttributes=title&data=file://home/mxb/code/srch/testdata/ndbooks.ndjson&attributesForFaceting=tags,authors,series&name=audiobooks`)
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

func TestIdxFindDocByPK(t *testing.T) {
	test := func(idx *Idx) error {
		docs, err := idx.findDocByPK(testDocPK)
		if err != nil {
			return err
		}
		var doc *Doc
		if len(docs) > 0 {
			doc = docs[0]
		}
		if doc.ID != testDocID {
			return fmt.Errorf("got %v doc id, wanted %v\n", doc.ID, testDocID)
		}
		if doc.PrimaryKey != testDocPK {
			return fmt.Errorf("got %v doc pk, wanted %v\n", doc.PrimaryKey, testDocPK)
		}
		return nil
	}
	runIdxTest(t, testIdxReq, test)
}

func TestIdxUpdateDoc(t *testing.T) {
	test := func(idx *Idx) error {
		r, err := idx.openData()
		if err != nil {
			return err
		}
		idx.getData = NdJSONFind(idx.PrimaryKey, r)
		d, err := idx.Find(testDocPK)
		if err != nil {
			return err
		}
		if len(d) < 1 {
			return fmt.Errorf("got %v results, expected at least one", len(d))
		}
		d[0]["title"] = "poot"

		err = idx.UpdateDoc(d[0])
		return nil
	}
	runIdxTest(t, testIdxReq, test)
}

func TestIdxFindData(t *testing.T) {
	test := func(idx *Idx) error {
		r, err := idx.openData()
		if err != nil {
			return err
		}
		idx.getData = NdJSONFind(idx.PrimaryKey, r)
		d, err := idx.Find(testDocPK)
		if err != nil {
			return err
		}
		id, ok := d[0][idx.PrimaryKey]
		if !ok {
			t.Errorf("data doesn't have pk, wanted %v\n", idx.PrimaryKey)
		}
		if pk := float64(testDocPK); id != pk {
			t.Errorf("got %v pk, wanted %v\n", id, pk)
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
