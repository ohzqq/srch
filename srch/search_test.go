package srch

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/ohzqq/cdb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
)

const (
	calLibPath = `/mnt/x/libraries/audiobooks`
)

type testSrchFunc func(*Idx, *Search) error

func runSrchTests(t *testing.T, test testSrchFunc) {
	for _, query := range TestQueryParams {
		runSrchTest(t, query, test)
	}
}

func runSrchTest(t *testing.T, query QueryStr, test testSrchFunc) {
	client, err := getTestClient(query)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := client.FindIdx(client.IndexName())
	if err != nil {
		t.Fatal(err)
	}
	test(idx, client.Search)
}

func calLib(ids []int) ([]map[string]any, error) {
	lib := cdb.NewLib(calLibPath)
	q := lib.NewQuery().GetByID(lo.ToAnySlice(resIDs)...)
	recs, err := lib.GetBooks(q)
	if err != nil {
		return nil, err
	}
	return recs.StringMap()
}

var resIDs = []int{1296, 1909, 2265, 2535, 2536, 2611, 2634, 4535, 4626, 5285, 5815, 5816, 6080, 6081, 6082, 6231, 6352, 6777, 6828, 6831, 6912, 7113}

func TestSearchRtrvAttr(t *testing.T) {
	test := func(idx *Idx, srch *Search) error {
		//wantResults, err := wantResults()
		u, err := url.Parse(ndjsonDataURL)
		if err != nil {
			return err
		}
		wantResults, err := FindData(u, resIDs)
		if err != nil {
			return err
		}

		if srch.Query != "" {
			gotResults, err := idx.Search(srch)
			if err != nil {
				t.Error(err)
			}
			for i, item := range gotResults {
				err = testTotalFields(srch.RtrvAttr, wantResults[i], item)
				if err != nil {
					t.Error(err)
				}
			}

		}
		return nil
	}
	runSrchTests(t, test)
}

func TestSearchDBData(t *testing.T) {
	test := func(idx *Idx, srch *Search) error {
		idx.SetFindDataFunc(calLib)

		wantResults, err := wantResults()
		if err != nil {
			return err
		}

		if srch.Query != "" {
			gotResults, err := idx.Search(srch)
			if err != nil {
				t.Error(err)
			}
			for i, item := range gotResults {
				want := cast.ToString(wantResults[i]["title"])
				got := cast.ToString(item["title"])
				if got != want {
					t.Errorf("got %v title, wanted %v\n", got, want)
				}
			}

		}
		return nil
	}
	runSrchTests(t, test)
}

func attrsToRtrv(attrs []string, data map[string]any) []string {
	var r []string
	for _, attr := range attrs {
		if attr == "*" {
			return maps.Keys(data)
		}
		r = append(r, attr)
	}
	return r
}

func testTotalFields(attr []string, test, res map[string]any) error {
	got := len(res)
	want := len(test)

	if t := len(attr); t > 0 {
		if attr[0] != "*" {
			want = t
		}
	}

	if got != want {
		return fmt.Errorf("got %v fields, wanted %v\n", got, want)
	}
	return nil
}

func wantResults() ([]map[string]any, error) {
	f, _ := os.Open(`/home/mxb/code/srch/testdata/ndbooks.ndjson`)
	defer f.Close()
	return findNDJSON(f, "id", resIDs)
}
