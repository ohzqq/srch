package srch

import (
	"net/url"
	"testing"
)

const testValuesCfg = `and=tags:count:desc&field=title&or=authors&or=narrators&or=series&data_dir=testdata/data-dir/`

func TestNewQuery(t *testing.T) {
	q := getNewQuery()
	if l := len(q); l != 6 {
		t.Errorf("got %v, expected %d\n,%v\n", l, 6, q)
	}
	i := New(testValuesCfg, testQueryString, testSearchString)
	if i.Len() != 7174 {
		t.Errorf("got %d, expected 7174\v", i.Len())
	}
	if len(i.Fields) != 5 {
		for _, f := range i.Fields {
			println(f.Attribute)
		}
		t.Errorf("got %d, expected %d\n", len(i.Fields), 5)
	}
	res := i.Search("fish")
	if len(res.Data) != 8 {
		t.Errorf("got %d, expected 8\n", len(res.Data))
	}
}

func getNewQuery() url.Values {
	return NewQuery(testValuesCfg, testQueryString, testSearchString)
}

func TestUrlValuesStringConfig(t *testing.T) {
	i, err := ParseCfgQuery(testValuesCfg)
	if err != nil {
		t.Error(err)
	}
	if len(i.Fields) != 5 {
		for _, f := range i.Fields {
			println(f.Attribute)
		}
		t.Errorf("got %d, expected %d\n", len(i.Fields), 5)
	}
}
