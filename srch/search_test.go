package srch

import (
	"fmt"
	"testing"
)

var bleveSearchTests = []string{
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`,
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`,
}

func TestBleveSearchAll(t *testing.T) {
	for i := 0; i < len(bleveSearchTests); i++ {
		q := bleveSearchTests[i]
		idx, err := New(q)
		if err != nil {
			t.Error(err)
		}
		if !idx.isBleve {
			t.Errorf("not bleve")
		}
		query := ""
		if i == 1 {
			query = "fish"
		}
		res, err := idx.Search(query)
		if err != nil {
			t.Error(err)
		}
		got := len(res)
		want := 7252
		if query == "fish" {
			want = 8
			//want = len(res)
		}
		err = searchErr(got, want, query)
		if err != nil {
			t.Error(err)
		}
	}
}

func searchErr(got int, want int, q string) error {
	if got != want {
		return fmt.Errorf("query %s got %d results, wanted %d\n", q, got, want)
	}
	return nil
}
