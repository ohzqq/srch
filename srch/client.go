package srch

import (
	"fmt"
	"slices"

	"github.com/ohzqq/hare"
	"github.com/samber/lo"
)

const (
	settingsTbl = "_settings"
	defaultTbl  = "default"
)

type Client struct {
	*Cfg

	db      *hare.Database
	tbl     *hare.Table
	indexes map[string]*Idx
}

func NewClient(cfg *Cfg) (*Client, error) {
	client := &Client{
		Cfg:     cfg,
		indexes: make(map[string]*Idx),
	}

	//step 1: initialize hare db
	err := client.initDB()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (client *Client) initDB() error {
	ds, err := NewDatastorage(client.DB())
	if err != nil {
		return fmt.Errorf("new datastorage error: %w\n", err)
	}

	err = client.SetDatastorage(ds)
	if err != nil {
		return fmt.Errorf("new set datastorage error: %w\n", err)
	}

	return nil
}

func (client *Client) LoadCfg() error {
	// check for settings table, create if it doesn't exist.
	if !client.db.TableExists(settingsTbl) {
		err := client.db.CreateTable(settingsTbl)
		if err != nil {
			return fmt.Errorf("client.LoadCfg error\n%w\n", err)
		}
		// since it doesn't exist, insert the current idx cfg
		_, err = client.db.Insert(settingsTbl, client.Idx)
		if err != nil {
			return err
		}
	}

	tbl, err := client.db.GetTable(settingsTbl)
	if err != nil {
		return err
	}
	client.tbl = tbl
	//client.SetTbl(tbl)

	ids, err := client.tbl.IDs()
	if err != nil {
		return err
	}

	for _, id := range ids {
		idx := &Idx{}
		err := client.tbl.Find(id, idx)
		if err != nil {
			return fmt.Errorf("%w: %v\n", err, client.IndexName())
		}
		client.indexes[idx.Name] = idx
	}
	return nil
}

func (client *Client) Indexes() map[string]*Idx {
	client.LoadCfg()
	return client.indexes
}

func (client *Client) HasIdx(name string) bool {
	idxs := client.Indexes()
	_, ok := idxs[name]
	return ok
}

func (client *Client) FindIdx(name string) (*Idx, error) {
	err := client.findIdxCfg(name)
	if err != nil {
		return nil, err
	}
	idxs := client.Indexes()
	idx, ok := idxs[name]
	if !ok {
		return nil, ErrIdxNotFound
	}
	return idx, nil
}

func (client *Client) FindIdxData(name string) (*Idx, error) {
	idx, err := client.FindIdx(name)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

func (client *Client) findIdxCfg(name string) error {

	// check if idx exists
	if !client.HasIdx(name) {
		// if it doesn't exist, create srch idx table
		if !client.db.TableExists(client.Idx.idxTblName()) {
			err := client.db.CreateTable(client.Idx.idxTblName())
			if err != nil {
				return err
			}

			// since the idx doesn't exist, insert param settings
			_, err = client.tbl.Insert(client.Idx)
			if err != nil {
				return err
			}
		}

		// create table for the index data
		if !client.db.TableExists(client.Idx.dataTblName()) {
			err := client.db.CreateTable(client.Idx.dataTblName())
			if err != nil {
				return err
			}
		}

		return nil
	}

	//rec, ok := client.indexes[name]
	//if !ok {
	//return errors.New("not ok")
	//}

	return nil
}

func (client *Client) SetDatastorage(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	client.db = h
	return nil
}

func (client *Client) TableNames() []string {
	client.LoadCfg()
	return lo.Without(client.db.TableNames(), settingsTbl, "")
}

func (client *Client) TableExists(name string) bool {
	return slices.Contains(client.TableNames(), name)
}
