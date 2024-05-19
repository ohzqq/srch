package param

import (
	"fmt"
	"slices"
	"testing"
)

func TestDecodeCfgStr(t *testing.T) {
	for query, test := range ParamTests() {
		cfg := NewCfg()
		err := cfg.Decode(query.String())
		if err != nil {
			t.Error(err)
		}

		IdxTests(t, query, cfg.Idx, test.Idx)
		CfgTests(t, query, cfg, test.Cfg)
		SrchTests(t, query, cfg.Search, test.Search)

	}
}

func TestDecodeCfgVals(t *testing.T) {
	for query, test := range ParamTests() {
		cfg := NewCfg()
		err := cfg.Decode(query.Query())
		if err != nil {
			t.Error(err)
		}

		IdxTests(t, query, cfg.Idx, test.Idx)
		CfgTests(t, query, cfg, test.Cfg)
		SrchTests(t, query, cfg.Search, test.Search)
	}
}

func sliceTest(num, field any, got, want []string) error {
	if !slices.Equal(got, want) {
		return paramTestMsg(num, field, got, want)
	}
	return nil
}

func paramTestMsg(num, field, got, want any) error {
	return fmt.Errorf("test %v, field %s\ngot %#v, wanted %#v\n", num, field, got, want)
}
