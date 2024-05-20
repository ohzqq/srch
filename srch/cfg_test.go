package srch

import (
	"errors"
	"testing"
)

type cfgTest struct {
	*Cfg
}

func TestDecodeCfgReq(t *testing.T) {
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Error(err)
		}
		test := req.cfgTest(i)

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
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.Query())
		if err != nil {
			t.Error(err)
		}
		test := req.cfgTest(i)

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
		test := req.cfgTest(i)

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

func (ct cfgTest) SrchCfg(got, want *Search) error {
	err := sliceErr("search.RtrvAttr", got.RtrvAttr, want.RtrvAttr)
	if err != nil {
		return err
	}
	err = sliceErr("search.Facets", got.Facets, want.Facets)
	if err != nil {
		return err
	}
	err = sliceErr("search.FacetFltr", got.FacetFltr, want.FacetFltr)
	if err != nil {
		return err
	}
	return nil
}

func (ct cfgTest) IdxCfg(got, want *IdxCfg) error {
	err := sliceErr("search.SrchAttr", got.SrchAttr, want.SrchAttr)
	if err != nil {
		return err
	}
	err = sliceErr("search.FacetAttr", got.FacetAttr, want.FacetAttr)
	if err != nil {
		return err
	}
	err = sliceErr("search.SortAttr", got.SortAttr, want.SortAttr)
	if err != nil {
		return err
	}
	return nil
}

func (ct cfgTest) cfg(got, want *Cfg) error {
	if got.IndexName() != want.IndexName() {
		return err(
			msg("cfg.IndexName()",
				got.IndexName(),
				want.IndexName(),
			),
			errors.New("index name doesn't match"),
		)
	}
	if got.Client.UID != want.Client.UID {
		return err(
			msg("cfg.Client.UID",
				got.Client.UID,
				want.Client.UID,
			),
			errors.New("index uid doesn't match"),
		)
	}
	if got.DataURL().Path != want.DataURL().Path {
		return err(
			msg("cfg.DataURL().Path",
				got.DataURL().Path,
				want.DataURL().Path,
			),
			errors.New("data path doesn't match"),
		)
	}
	if got.DB().Path != want.DB().Path {
		return err(
			msg("cfg.DB().Path",
				got.DB().Path,
				want.DB().Path,
			),
			errors.New("db path doesn't match"),
		)
	}
	if got.SrchURL().Path != want.SrchURL().Path {
		return err(
			msg("cfg.SrchURL().Path",
				got.SrchURL().Path,
				want.SrchURL().Path),
			errors.New("srch path doesn't match"),
		)
	}
	return nil
}
