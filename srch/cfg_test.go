package srch

import (
	"errors"
	"fmt"
	"testing"
)

type cfgTest struct {
	*Cfg
}

func TestNewClient(t *testing.T) {
	for _, query := range TestQueryParams {
		cfg, err := NewCfg(query.Query())
		if err != nil {
			t.Error(err)
		}
		client, err := NewClient(cfg)
		if err != nil {
			t.Error(err)
		}
		_, err = client.FindIdx(client.IndexName())
		if err != nil {
			t.Errorf("name %v\n%v\n", query.Query(), err)
		}
	}

	runTests(t, testHasTbls)
}

func testHasTbls(_ int, req reqTest) error {
	client, err := req.Client()
	if err != nil {
		return err
	}
	idx, err := client.FindIdx(client.IndexName())
	if err != nil {
		return err
	}
	if got := !client.db.TableExists(idx.idxTblName()); got {
		want := true
		if got != want {
			return fmt.Errorf("got %v for tbl %v, wanted %v\n", got, idx.idxTblName(), want)
		}
	}
	if got := !client.db.TableExists(idx.dataTblName()); got {
		want := true
		if got != want {
			return fmt.Errorf("got %v for tbl %v, wanted %v\n", got, idx.dataTblName(), want)
		}
	}

	return nil
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

func TestCfgChanged(t *testing.T) {
	var wanted = []bool{
		false,
		true,
		true,
		true,
		true,
	}
	for i, query := range changedCfg {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Error(err)
		}
		//want := req.cfgWant(i)
		cfg, err := req.cfgGot()
		if err != nil {
			t.Error(err)
		}

		r, err := NewRequest(TestQueryParams[len(TestQueryParams)-1].String())
		if err != nil {
			t.Error(err)
		}
		test, err := r.Cfg()
		if err != nil {
			t.Error(err)
		}

		got := test.Idx.Changed(cfg.Idx)
		want := wanted[i]

		if got != want {
			t.Errorf("query %v\ncfg change %v, %v\nparam %#v\ndb %#v\n", query.String(), cfg.IndexName(), test.Idx.Changed(cfg.Idx), test.Idx, cfg.Idx)
		}
	}
}

func TestUpdateChangedCfg(t *testing.T) {
	for _, query := range changedCfg {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Error(err)
		}
		client, err := req.Client()
		if err != nil {
			t.Error(err)
		}
		idx, err := client.FindIdx(client.IndexName())
		if err != nil {
			t.Error(err)
		}

		if idx.Changed(client.Idx) {
			client.Idx.SetID(idx.GetID())
			err := client.tbl.Update(client.Idx)
			if err != nil {
				t.Error(err)
			}
			//println(client.Idx.GetID())
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

func (ct cfgTest) IdxCfg(got, want *Idx) error {
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
	wm := want.mapParams()
	for k, v := range got.Mapping {
		if _, ok := wm[k]; !ok {
			return errors.New("no key")
		}
		err := strSliceErr("doc.Mapping", v, wm[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func (ct cfgTest) cfg(got, want *Cfg) error {
	if got.IndexName() != want.IndexName() {
		return newErr(
			msg("cfg.IndexName()",
				got.IndexName(),
				want.IndexName(),
			),
			errors.New("index name doesn't match"),
		)
	}
	if got.Idx.UID != want.Idx.UID {
		return newErr(
			msg("cfg.Client.UID",
				got.Idx.UID,
				want.Idx.UID,
			),
			errors.New("index uid doesn't match"),
		)
	}
	if got.DataURL().Path != want.DataURL().Path {
		return newErr(
			msg("cfg.DataURL().Path",
				got.DataURL().Path,
				want.DataURL().Path,
			),
			errors.New("data path doesn't match"),
		)
	}
	if got.DB().Path != want.DB().Path {
		return newErr(
			msg("cfg.DB().Path",
				got.DB().Path,
				want.DB().Path,
			),
			errors.New("db path doesn't match"),
		)
	}
	if got.SrchURL().Path != want.SrchURL().Path {
		return newErr(
			msg("cfg.SrchURL().Path",
				got.SrchURL().Path,
				want.SrchURL().Path),
			errors.New("srch path doesn't match"),
		)
	}
	return nil
}
