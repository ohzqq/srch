package param

import (
	"fmt"
	"slices"
	"testing"
)

func TestDecodeCfgStr(t *testing.T) {
	for query, test := range cfgTests {
		cfg := NewCfg()
		err := cfg.Decode(query.String())
		if err != nil {
			t.Error(err)
		}

		testIdx(t, query, cfg.Idx, test.Idx)
		testCfg(t, query, cfg, test.Cfg)
		testSrch(t, query, cfg.Search, test.Search)

	}
}

func TestDecodeCfgVals(t *testing.T) {
	for query, test := range cfgTests {
		cfg := NewCfg()
		err := cfg.Decode(query.Query())
		if err != nil {
			t.Error(err)
		}

		testIdx(t, query, cfg.Idx, test.Idx)

		if cfg.IndexName() != test.IndexName() {
			t.Errorf("test %v Index: got %#v, expected %#v\n", query, cfg.IndexName(), test.IndexName())
		}
		if cfg.Client.UID != test.Client.UID {
			t.Errorf("test %v ID: got %#v, expected %#v\n", query, cfg.Client.UID, test.Client.UID)
		}
		if cfg.DataURL().Path != test.DataURL().Path {
			t.Errorf("test %v\ndata Path:\ngot %#v, expected %#v\n", query, cfg.DataURL().Path, test.DataURL().Path)
		}
		if cfg.DB().Path != test.DB().Path {
			t.Errorf("test %v\ndb Path:\ngot %#v, expected %#v\n", query, cfg.DB().Path, test.DB().Path)
		}
		if cfg.SrchURL().Path != test.SrchURL().Path {
			t.Errorf("test %v\nsrch Path:\ngot %#v, expected %#v\n", query, cfg.SrchURL().Path, test.SrchURL().Path)
		}
	}
}

func TestEncodeCfg(t *testing.T) {
	t.SkipNow()
	for num, test := range cfgTests {
		v, err := Encode(test.Idx)
		if err != nil {
			t.Error(err)
		}
		if v.Encode() != test.Query {
			t.Errorf("test %v: got %v, wanted %v\n", num, v.Encode(), test.Query)
		}
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
