package srch

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
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

var resIDs = []int{1296, 1909, 2265, 2535, 2536, 2611, 2634, 4535, 4626, 5285, 5815, 5816, 6080, 6081, 6082, 6231, 6352, 6777, 6828, 6831, 6912, 7113}

func TestSearch(t *testing.T) {
	test := func(idx *Idx, srch *Search) error {
		r, err := idx.openData()
		if err != nil {
			return err
		}
		idx.getData = NdJSONFind(idx.PrimaryKey, r)

		wantResults, err := wantResults()
		if err != nil {
			return err
		}

		if srch.Query != "" {
			gotResults, err := idx.Search(srch)
			if err != nil {
				t.Error(err)
			}
			titles := make([]string, len(gotResults))
			for i, item := range gotResults {
				titles[i] = cast.ToString(item["title"])
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
		println(attr[0] == "*")
	}
	if got != want {
		return fmt.Errorf("got %v fields, wanted %v\n", got, want)
	}
	return nil
}

func wantResults() ([]map[string]any, error) {
	f, _ := os.Open(`/home/mxb/code/srch/testdata/ndbooks.ndjson`)
	defer f.Close()
	return findNDJSON(f, "id", resIDs...)
}
