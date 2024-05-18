package index

import (
	"errors"
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/hare/dberr"
	"github.com/ohzqq/srch/param"
)

const (
	settingsTbl = "_settings"
	defaultTbl  = "default"
)

type Client struct {
	*hare.Database
	*param.Client

	tbl *hare.Table
}

func New(settings any) (*Client, error) {
	client := &Client{
		Client: param.NewClient(),
	}

	err := param.Decode(settings, client.Client)
	if err != nil {
		return nil, fmt.Errorf("param decoding error: %w\n", err)
	}

	ds, err := NewDatastorage(client.Client.URL)
	if err != nil {
		return nil, fmt.Errorf("new datastorage error: %w\n", err)
	}

	err = client.SetDatastorage(ds)
	if err != nil {
		return nil, fmt.Errorf("new set datastorage error: %w\n", err)
	}

	err = client.init()
	if err != nil {
		return nil, fmt.Errorf("client init error: %w\n", err)
	}

	return client, nil
}

func (client *Client) init() error {
	//Create settings table if it doesn't exist
	if !client.Database.TableExists(settingsTbl) {
		err := client.Database.CreateTable(settingsTbl)
		if err != nil {
			return fmt.Errorf("db.getCfg CreateTable error\n%w\n", err)
		}

		//cfg, err := client.Cfg()
		//if err != nil {
		//return err
		//}
		//err = cfg.Insert(NewCfgParams(client.Params))
		//if err != nil {
		//return fmt.Errorf("client.init cfg.Insert error\n%w\n", err)
		//}
	}

	return nil
}

func (client *Client) Cfg() (*Client, error) {
	tbl, err := client.Database.GetTable(settingsTbl)
	if err != nil {
		return nil, err
	}
	client.SetTbl(tbl)
	return client, nil
}

func (client *Client) GetIdxCfg(name string) (*IdxCfg, error) {
	return client.findIdxCfg(name, "cfg")
}

func (client *Client) GetIdx(name string) (*Idx, error) {
	idx, err := client.findIdxCfg(name, "table")
	if err != nil {
		return nil, err
	}
	return NewIdx(client.Database, idx), nil
}

func (client *Client) findIdxCfg(name, create string) (*IdxCfg, error) {
	clientCfg, err := client.Cfg()
	if err != nil {
		return nil, err
	}

	//get idx cfg from params
	//cur := NewCfgParams(client.Params)
	cur := NewCfg()

	//find existing cfg
	idxCfg, err := clientCfg.Find(name)
	if err == nil {
		//if cfg is found, set cur.ID to cfg.ID
		cur.SetID(idxCfg.GetID())
	} else if err != nil {
		switch {
		case errors.Is(err, dberr.ErrNoTable):
			//when getting the index, create table if it doesn't exist or insert cfg
			//for the table if it doesn't exist.
			switch create {
			case "table":
				err = client.Database.CreateTable(clientCfg.Index)
			case "cfg":
				//if inserting a new cfg to settings, set idxCfg to cur so that it isn't
				//nil
				idxCfg = cur
				err = clientCfg.Insert(cur)
			}
			if err != nil && !errors.Is(err, dberr.ErrTableExists) {
				return nil, fmt.Errorf("client.GetIdxCfg create table error:\n%w: %v\n", err, clientCfg.Index)
			}
		default:
			return nil, fmt.Errorf("client.GetIdxCfg error:\n%w: %v\n", err, clientCfg.Index)
		}
	}

	//check to see if the provided client.Params are different from the database
	//record, if so, update.
	if !param.CfgEqual(idxCfg.Idx, cur.Idx) {
		err := clientCfg.Update(cur)
		if err != nil {
			return nil, err
		}
	}

	return idxCfg, nil
}

func (cfg *Client) Insert(idx *IdxCfg) error {
	_, err := cfg.tbl.Insert(idx)
	if err != nil {
		return fmt.Errorf("cfg.Insert error\n%w\n", err)
	}
	return nil
}

func (cfg *Client) Update(idx *IdxCfg) error {
	err := cfg.tbl.Update(idx)
	if err != nil {
		return fmt.Errorf("cfg.Insert error\n%w\n", err)
	}
	return nil
}

func (cfg *Client) Find(name string) (*IdxCfg, error) {
	ids, err := cfg.tbl.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		idx := &IdxCfg{}
		err := cfg.tbl.Find(id, idx)
		if err != nil {
			return nil, err
		}
		if idx.Index == name {
			return idx, nil
		}
	}

	return nil, dberr.ErrNoTable
}

func (cfg *Client) Tables() ([]*IdxCfg, error) {
	ids, err := cfg.tbl.IDs()
	if err != nil {
		return nil, err
	}

	tbls := make([]*IdxCfg, len(ids))

	for i, id := range ids {
		idx := &IdxCfg{}
		err := cfg.tbl.Find(id, idx)
		if err != nil {
			return nil, err
		}
		tbls[i] = idx
	}
	return tbls, nil
}

func (cfg *Client) SetTbl(tbl *hare.Table) *Client {
	cfg.tbl = tbl
	return cfg
}

func (client *Client) SetDatastorage(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	client.Database = h
	return nil
}

func (client *Client) memDB() error {
	r := &ram.Ram{
		Store: store.New(),
	}
	return client.SetDatastorage(r)
}
