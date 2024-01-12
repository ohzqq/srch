package srch

import (
	"net/url"
	"testing"
)

const testValuesCfg = `and=tags&field=title&or=authors&or=narrators&or=series&file=testdata/config.json`

func TestNewQuery(t *testing.T) {
	q := getNewQuery()
	if l := len(q); l != 5 {
		t.Errorf("got %v, expected %d\n,%v\n", l, 5, q)
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
