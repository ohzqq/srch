package srch

import "testing"

var bleveSearchTests = []string{
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id`,
	`searchableAttributes=title&fullText=../testdata/poot.bleve&uid=id&query=fish`,
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
	}
}
