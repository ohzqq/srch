package srch

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/exp/maps"
)

type clientTest struct {
	*Client
	got  *Client
	want *Client
}

func TestClientMem(t *testing.T) {
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Fatal(err)
		}
		test, err := req.clientTest(i)

		err = test.storage(query)
		if err != nil {
			t.Error(err)
		}

		err = test.settingsExists(query)
		if err != nil {
			t.Error(err)
		}

		err = test.getClientCfg(query)
		if err != nil {
			t.Error(err)
		}

		err = test.getIdxCfg(query)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestClientTotalIdxs(t *testing.T) {
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Fatal(err)
		}
		test, err := req.clientTest(i)
		err = test.testTotalIdx(query)
		if err != nil {
			t.Error(err)
		}
	}

}

func TestClientMemListTbls(t *testing.T) {
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Fatal(err)
		}
		test, err := req.clientTest(i)

		err = test.listTbls(query)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFindIdx(t *testing.T) {
	for i, query := range TestQueryParams {
		req, err := newTestReq(query.String())
		if err != nil {
			t.Fatal(err)
		}
		test, err := req.clientTest(i)
		if err != nil {
			t.Error(err)
		}

		_, err = test.got.FindIdx(test.got.IndexName())
		if err != nil {
			t.Error(err)
		}
		_, err = test.want.FindIdx(test.want.IndexName())
		if err != nil {
			t.Error(err)
		}

	}
}

func TestClientDisk(t *testing.T) {
}

func TestClientNet(t *testing.T) {
}

func (t clientTest) storage(q QueryStr) error {
	gs, err := NewDatastorage(t.got.DB())
	if err != nil {
		return fmt.Errorf("new datastorage error: %w\n", err)
	}
	ws, err := NewDatastorage(t.want.DB())
	if err != nil {
		return fmt.Errorf("new datastorage error: %w\n", err)
	}
	if reflect.TypeOf(gs).String() != reflect.TypeOf(ws).String() {
		return errors.New(msg(q.String(), gs, ws))
	}
	return nil
}

func (t clientTest) settingsExists(q QueryStr) error {
	got := t.got.db.TableExists(settingsTbl)
	want := t.want.db.TableExists(settingsTbl)
	if got != want {
		return errors.New(msg(q.String(), got, want))
	}
	return nil
}

func (t clientTest) getClientCfg(q QueryStr) error {
	err := t.got.LoadCfg()
	err = t.want.LoadCfg()
	terr := errors.New(msg(q.String(), t.got, t.want))
	if err != nil {
		return terr
	}

	err = t.testIDs(q, terr)
	if err != nil {
		return err
	}

	return nil
}

func (t clientTest) getIdxCfg(q QueryStr) error {
	gn := t.got.IndexName()
	//println(gn)
	gCfg, err := t.got.FindIdx(gn)
	if err != nil {
		return errors.New(msg(q.String(), gCfg, t.got))
	}

	wn := t.want.IndexName()
	wCfg, err := t.want.FindIdx(wn)
	if err != nil {
		return err
	}
	if gCfg.ID != wCfg.ID {
		return errors.New("not same id")
	}

	return nil
}

func (t clientTest) testIDs(q QueryStr, terr error) error {
	err := t.got.LoadCfg()
	if err != nil {
		return terr
	}
	err = t.want.LoadCfg()
	if err != nil {
		return terr
	}
	gIDs, err := t.got.tbl.IDs()
	if err != nil {
		return terr
	}
	wIDs, err := t.want.tbl.IDs()
	if err != nil {
		return terr
	}

	return intSliceErr(q.String(), gIDs, wIDs)
}

func (t clientTest) testTotalIdx(q QueryStr) error {

	gCfgs := t.got.Indexes()
	//if err != nil {
	//  return nil
	//}

	gIDs, err := t.got.tbl.IDs()
	if err != nil {
		return err
	}

	wCfgs := t.want.Indexes()
	//if err != nil {
	//  return err
	//}

	if len(gCfgs) != len(wCfgs) {
		return errors.New("len of configs not equal")
	}

	if len(gIDs) != len(gCfgs) {
		return fmt.Errorf("got %v _settings records, wanted %v\n", len(gIDs), len(gCfgs))
	}

	return nil
}

func (t clientTest) listTbls(q QueryStr) error {
	got := maps.Keys(t.got.Indexes())
	want := maps.Keys(t.want.Indexes())
	if len(want) == 0 {
		want = []string{"default"}
	}
	//want := []string{t.want.IndexName()}
	return strSliceErr(q.String(), got, want)
}
