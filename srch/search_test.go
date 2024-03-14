package srch

import "testing"

var bleveSearchTests = []string{
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`,
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`,
}

func TestBleveSearchAll(t *testing.T) {
	println("what??")
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
		println(len(res))
	}
}
