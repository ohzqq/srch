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
	Params *param.Cfg
}

func New(settings any) (*Client, error) {
	client := &Client{
		Params: param.NewCfg(),
	}

	err := param.Decode(settings, client.Params)
	if err != nil {
		return nil, fmt.Errorf("param decoding error: %w\n", err)
	}

	ds, err := NewDatastorage(client.Params.URL)
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

		cfg, err := client.Cfg()
		if err != nil {
			return err
		}
		err = cfg.Insert(NewCfgParams(client.Params))
		if err != nil {
			return fmt.Errorf("client.init cfg.Insert error\n%w\n", err)
		}
	}

	return nil
}

func (client *Client) Cfg() (*ClientCfg, error) {
	cfg := NewClientCfg(client.Params)
	tbl, err := client.Database.GetTable(settingsTbl)
	if err != nil {
		return nil, err
	}
	cfg.SetTbl(tbl)
	return cfg, nil
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
	cfg, err := client.Cfg()
	if err != nil {
		return nil, err
	}
	idx, err := cfg.Find(name)
	if err != nil {
		switch {
		case errors.Is(err, dberr.ErrNoTable):
			var err error
			switch create {
			case "table":
				err = client.Database.CreateTable(cfg.Index)
			case "cfg":
				err = cfg.Insert(NewCfgParams(client.Params))
			}
			if err != nil && !errors.Is(err, dberr.ErrTableExists) {
				return nil, fmt.Errorf("client.GetIdxCfg create table error:\n%w: %v\n", err, cfg.Index)
			}
		default:
			return nil, fmt.Errorf("client.GetIdxCfg error:\n%w: %v\n", err, cfg.Index)
		}
	}
	return idx, nil
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
