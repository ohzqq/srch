package srch

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
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

func (t clientTest) listTbls(q QueryStr) error {
	got := t.got.Database.TableNames()
	want := t.want.Database.TableNames()
	slices.Sort(got)
	slices.Sort(want)
	return sliceErr(q.String(), got, want)
}
