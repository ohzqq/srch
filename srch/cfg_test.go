package srch

import (
	"testing"
)

func TestDecodeCfgStr(t *testing.T) {
	for query, test := range ParamTests() {
		cfg := NewCfg()
		err := cfg.Decode(query.Query())
		if err != nil {
			t.Error(err)
		}

		err = test.Srch(cfg.Search, test.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = test.Index(cfg.Idx, test.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = test.Config(cfg, test.Cfg)
		if err != nil {
			t.Errorf("cfg test query %v\n%#v\n", query.String(), err)
		}
	}
}

func TestDecodeCfgVals(t *testing.T) {
	for query, test := range ParamTests() {
		cfg := NewCfg()
		err := cfg.Decode(query.Query())
		if err != nil {
			t.Error(err)
		}

		err = test.Srch(cfg.Search, test.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = test.Index(cfg.Idx, test.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = test.Config(cfg, test.Cfg)
		if err != nil {
			t.Errorf("cfg test query %v\n%#v\n", query.String(), err)
		}
	}
}
