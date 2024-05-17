package index

import (
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
		err = client.SetCfg(DefaultCfg())
		if err != nil {
			return fmt.Errorf("db.getCfg Insert error\n%w\n", err)
		}
	}
	return nil
}

func (client *Client) GetIdx(name string) (*Idx, error) {
	cfg, err := client.GetCfg(name)
	if err != nil {
		return nil, err
	}
	return NewIdx(client.Database, cfg), nil
}

func (client *Client) SetCfg(cfg *Cfg) error {
	_, err := client.Database.Insert(settingsTbl, cfg)
	if err != nil {
		return fmt.Errorf("db.getCfg Insert error\n%w\n", err)
	}
	return nil
}

func (client *Client) GetCfg(name string) (*Cfg, error) {
	tbl, err := client.Cfg()
	if err != nil {
		return nil, err
	}

	ids, err := tbl.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		cfg := &Cfg{}
		err := client.Database.Find(settingsTbl, id, cfg)
		if err != nil {
			return nil, err
		}
		if cfg.Index == name {
			return cfg, nil
		}
	}

	return nil, dberr.ErrNoTable
}

func (client *Client) Cfg() (*hare.Table, error) {
	return client.Database.GetTable(settingsTbl)
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
