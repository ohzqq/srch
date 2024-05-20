package srch

import (
	"testing"
)

func TestDecodeCfgReq(t *testing.T) {
	for _, query := range TestQueryParams {
		test := getCfgParamTest(query)
		req, err := newTestReq(query.String())
		if err != nil {
			t.Error(err)
		}

		cfg, err := req.Cfg()
		if err != nil {
			t.Error(err)
		}

		err = test.SrchCfg(cfg.Search, test.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = test.IdxCfg(cfg.Idx, test.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = test.cfg(cfg, test.Cfg)
		if err != nil {
			t.Errorf("cfg test query %v\n%#v\n", query.String(), err)
		}
	}
}

func TestDecodeCfgVals(t *testing.T) {
	for _, query := range TestQueryParams {
		test := getCfgParamTest(query)
		req, err := newTestReq(query.Query())
		if err != nil {
			t.Error(err)
		}

		cfg, err := req.Cfg()
		if err != nil {
			t.Error(err)
		}

		err = test.SrchCfg(cfg.Search, test.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = test.IdxCfg(cfg.Idx, test.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = test.cfg(cfg, test.Cfg)
		if err != nil {
			t.Errorf("cfg test query %v\n%#v\n", query.String(), err)
		}
	}
}

func TestDecodeCfgStr(t *testing.T) {
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Error(err)
		}
		test := req.cfgTest(getTestCfg(i))

		cfg, err := req.Cfg()
		if err != nil {
			t.Error(err)
		}

		err = test.SrchCfg(cfg.Search, test.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = test.IdxCfg(cfg.Idx, test.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = test.cfg(cfg, test.Cfg)
		if err != nil {
			t.Errorf("cfg test query %v\n%#v\n", query.String(), err)
		}
	}
}
