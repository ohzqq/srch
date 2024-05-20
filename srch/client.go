package srch

import (
	"fmt"

	"github.com/ohzqq/hare"
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

func (client *Client) SetDatastorage(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	client.Database = h
	return nil
}
