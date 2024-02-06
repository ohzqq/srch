package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7174

const testDataFile = `testdata/data-dir/audiobooks.json`
const testDataDir = `testdata/data-dir`
const testCfgFile = `testdata/config-old.json`
const testYAMLCfgFile = `testdata/config.yaml`
const testCfgFileData = `testdata/config-with-data.json`
const libCfgStr = "searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&dataFile=testdata/data-dir/audiobooks.json"

func TestData(t *testing.T) {
	books = loadData(t)
	err := totalBooksErr(len(books), 71734)
	if err != nil {
		t.Error(err)
	}
}

var testQueryNewIndex = []string{
	"searchableAttributes=title&fullText",
	"searchableAttributes=title&dataDir=testdata/data-dir",
	"attributesForFaceting=tags,authors,series,narrators&dataFile=testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&dataFile=testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators",
}

var titleField = NewField(DefaultField)

func TestNewIndex(t *testing.T) {
	for i := 0; i < len(testQuerySettings); i++ {
		q := testQuerySettings[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		if !idx.HasData() {
			data := loadData(t)
			idx.Index(data)
		}
		err = totalBooksErr(idx.Len(), q)
		if err != nil {
			t.Error(err)
		}
	}
}

func intErr(got, want int, msg ...any) error {
	if got != want {
		err := fmt.Errorf("got %d, want %d\n", got, want)
		if len(msg) > 0 {
			err = fmt.Errorf("%w [msg] %v\n", err, msg)
		}
		return err
	}
	return nil
}

func totalBooksTest(total int, t *testing.T) {
	if total != 7174 {
		t.Errorf("got %d, expected %d\n", total, 7174)
	}
}

func newTestIdx() *Index {
	//idx, err := New(libCfgStr)
	//if err != nil {
	//log.Fatal(err)
	//}
	return newTestIdxCfg("")
}

func newTestIdxCfg(p string) *Index {
	idx, err := New(libCfgStr + p)
	if err != nil {
		log.Fatal(err)
	}
	return idx
}

func TestSortIndexByTitle(t *testing.T) {
	title := libCfgStr + "&sortBy=title"
	idx, err := New(title)
	if err != nil {
		t.Error(err)
	}
	title, ok := idx.Data[0][DefaultField].(string)
	if !ok {
		t.Errorf("not a string")
	}
	if title != "#Blur" {
		t.Errorf("sorting err, got %s, expected %s\n", title, "#Blur")
	}
	//fmt.Printf("%+v\n", idx.Data[0])
}

func TestSortIndexByDate(t *testing.T) {
	title := libCfgStr + "&sortableAttributes=title:text&sortableAttributes=added_stamp:num" + "&sortBy=added_stamp&order=desc"
	idx, err := New(title)
	if err != nil {
		log.Fatal(err)
	}
	title, ok := idx.Data[0][DefaultField].(string)
	if !ok {
		t.Errorf("not a string")
	}
	if title != "Risk the Fall" {
		t.Errorf("sorting err, got %s, expected %s\n", title, "Risk the Fall")
	}
}

func totalBooksErr(total int, vals ...any) error {
	if total != 7174 {
		err := fmt.Errorf("got %d, expected %d\n", total, 7174)
		return fmt.Errorf("%w\nmsg: %v", err, vals)
	}
	return nil
}

func loadData(t *testing.T) []map[string]any {
	d, err := os.ReadFile(testDataFile)
	if err != nil {
		t.Error(err)
	}

	var books []map[string]any
	err = json.Unmarshal(d, &books)
	if err != nil {
		t.Error(err)
	}

	books = books

	return books
}
