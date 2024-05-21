package srch

import (
	"errors"
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/dberr"
)

const (
	settingsTbl = "_settings"
	defaultTbl  = "default"
)

type Client struct {
	*Cfg
	*hare.Database

	tbl *hare.Table
}

func NewClient(cfg *Cfg) (*Client, error) {
	client := &Client{
		Cfg: cfg,
	}

	ds, err := NewDatastorage(client.DB())
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
			return fmt.Errorf("client.init error\n%w\n", err)
		}
	}
	return nil
}

func (client *Client) getCfgTbl() error {
	tbl, err := client.Database.GetTable(settingsTbl)
	if err != nil {
		return err
	}
	client.SetTbl(tbl)
	return nil
}

func (client *Client) IdxIDs() ([]int, error) {
	return client.tbl.IDs()
}

func (client *Client) FindIdxCfg(name string) (*IdxCfg, error) {
	ids, err := client.tbl.IDs()
	if err != nil {
		return nil, fmt.Errorf("%w: %v\n", err, client.IndexName())
	}

	for _, id := range ids {
		idx := &IdxCfg{}
		err := client.tbl.Find(id, idx)
		if err != nil {
			return nil, fmt.Errorf("%w: %v\n", err, client.IndexName())
		}
		if client.IndexName() == name {
			return idx, nil
		}
	}

	return nil, dberr.ErrNoTable
}

func (client *Client) findIdxCfg(name, create string) (*IdxCfg, error) {
	err := client.getCfgTbl()
	if err != nil {
		return nil, err
	}

	//get idx cfg from params
	//cur := NewCfgParams(client.Params)
	cur := NewIdxCfg()

	//find existing cfg
	idxCfg, err := client.FindIdxCfg(name)
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
				err = client.Database.CreateTable(client.IndexName())
			case "cfg":
				//if inserting a new cfg to settings, set idxCfg to cur so that it isn't
				//nil
				idxCfg = cur
				_, err = client.tbl.Insert(cur)
			}
			if err != nil && !errors.Is(err, dberr.ErrTableExists) {
				return nil, fmt.Errorf("client.GetIdxCfg create table error:\n%w: %v\n", err, client.IndexName())
			}
		default:
			return nil, fmt.Errorf("client.GetIdxCfg error:\n%w: %v\n", err, client.IndexName())
		}
	}

	//check to see if the provided client.Params are different from the database
	//record, if so, update.
	//if !param.CfgEqual(idxCfg.Idx, cur.Idx) {
	//  err := clientCfg.Update(cur)
	//  if err != nil {
	//    return nil, err
	//  }
	//}

	return idxCfg, nil
}

func (client *Client) SetDatastorage(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	client.Database = h
	return nil
}

func (cfg *Client) SetTbl(tbl *hare.Table) *Client {
	cfg.tbl = tbl
	return cfg
}
