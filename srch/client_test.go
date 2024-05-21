package srch

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/samber/lo"
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
	got := t.got.Database.TableExists(settingsTbl)
	want := t.want.Database.TableExists(settingsTbl)
	if got != want {
		return errors.New(msg(q.String(), got, want))
	}
	return nil
}

func (t clientTest) getClientCfg(q QueryStr) error {
	err := t.got.GetCfg()
	err = t.want.GetCfg()
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
	gCfg, err := t.got.FindIdxCfg(gn)
	if err != nil {
		return err
	}

	wn := t.want.IndexName()
	wCfg, err := t.want.FindIdxCfg(wn)
	if err != nil {
		return err
	}
	if gCfg.ID != wCfg.ID {
		return errors.New("not same id")
	}

	return nil
}

func (t clientTest) testIDs(q QueryStr, terr error) error {
	gIDs, err := t.got.IdxIDs()
	if err != nil {
		return terr
	}
	wIDs, err := t.want.IdxIDs()
	if err != nil {
		return terr
	}

	return intSliceErr(q.String(), gIDs, wIDs)
}

func (t clientTest) listTbls(q QueryStr) error {
	got := t.got.TableNames()
	want := lo.Without(t.want.Database.TableNames(), "", "_settings")
	if len(want) == 0 {
		want = []string{"default"}
	}
	//want := []string{t.want.IndexName()}
	return strSliceErr(q.String(), got, want)
}
