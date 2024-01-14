package srch

import (
	"fmt"
	"net/url"
	"testing"
)

const testValuesCfg = `and=tags:count:desc&field=title&or=authors:label:asc&or=narrators&or=series&data_dir=testdata/data-dir/`

func TestNewQuery(t *testing.T) {
	//t.SkipNow()
	q := getNewQuery()
	if l := len(q); l != 6 {
		t.Errorf("got %v, expected %d\n,%v\n", l, 6, q)
	}
	query := fmt.Sprintf("%s&%s&%s", testValuesCfg, testQueryString, testSearchString)
	i := New(query)
	if i.Len() != 3 {
		t.Errorf("got %d, expected %d\v", i.Len(), 3)
	}
	if len(i.Fields) != 5 {
		for _, f := range i.Fields {
			println(f.Attribute)
		}
		t.Errorf("got %d, expected %d\n", len(i.Fields), 5)
	}
	//res, err := json.Marshal(i)
	//if err != nil {
	//  t.Error(err)
	//}

	//i.PrettyPrint()

	//println(string(res))
	//n := &Index{}
	//err = json.Unmarshal(res, n)
	//if err != nil {
	//t.Error(err)
	//}
	//fmt.Printf("%v\n", n)
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
