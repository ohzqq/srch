package srch

import (
	"errors"
	"fmt"
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

		want := req.cfgWant(i)
		got, err := req.cfgGot()
		if err != nil {
			t.Error(err)
		}

		err = want.SrchCfg(got.Search, want.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = want.IdxCfg(got.Idx, want.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = want.cfg(got, want.Cfg)
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
		want := req.cfgWant(i)
		got, err := req.cfgGot()
		if err != nil {
			t.Error(err)
		}

		err = want.SrchCfg(got.Search, want.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = want.IdxCfg(got.Idx, want.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = want.cfg(got, want.Cfg)
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
		want := req.cfgWant(i)
		got, err := req.cfgGot()
		if err != nil {
			t.Error(err)
		}

		err = want.SrchCfg(got.Search, want.Search)
		if err != nil {
			t.Errorf("srch test query %v\n%#v\n", query.String(), err)
		}

		err = want.IdxCfg(got.Idx, want.Idx)
		if err != nil {
			t.Errorf("idx test query %v\n%#v\n", query.String(), err)
		}

		err = want.cfg(got, want.Cfg)
		if err != nil {
			t.Errorf("cfg test query %v\n%#v\n", query.String(), err)
		}
	}
}

func (ct cfgTest) SrchCfg(got, want *Search) error {
	err := strSliceErr("search.RtrvAttr", got.RtrvAttr, want.RtrvAttr)
	if err != nil {
		return err
	}
	err = strSliceErr("search.Facets", got.Facets, want.Facets)
	if err != nil {
		return err
	}
	err = strSliceErr("search.FacetFltr", got.FacetFltr, want.FacetFltr)
	if err != nil {
		return err
	}
	return nil
}

func (ct cfgTest) IdxCfg(got, want *IdxCfg) error {
	err := strSliceErr("search.SrchAttr", got.SrchAttr, want.SrchAttr)
	if err != nil {
		return err
	}
	err = strSliceErr("search.FacetAttr", got.FacetAttr, want.FacetAttr)
	if err != nil {
		return err
	}
	err = strSliceErr("search.SortAttr", got.SortAttr, want.SortAttr)
	if err != nil {
		return err
	}
	fmt.Printf("got map %#v\nwant map %#v\n", got.Mapping, want.Mapping)
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
