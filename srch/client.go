package srch

import (
	"fmt"
	"net/url"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
)

const (
	settingsTbl = "_settings"
	defaultTbl  = "default"
)

type Client struct {
	*hare.Database
	*url.URL `json:"-"`
	*ClientCfg

	tbl *hare.Table
}

func NewClient(settings any) (*Client, error) {
	client := &Client{
		URL: &url.URL{},
	}

	err := Decode(settings, client)
	if err != nil {
		return nil, fmt.Errorf("param decoding error: %w\n", err)
	}

	ds, err := NewDatastorage(client.URL)
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

func (client *ClientCfg) init() error {
	//Create settings table if it doesn't exist
	if !client.Database.TableExists(settingsTbl) {
		err := client.Database.CreateTable(settingsTbl)
		if err != nil {
			return fmt.Errorf("client.init error\n%w\n", err)
		}
	}
	return nil
}

func (client *ClientCfg) SetDatastorage(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	client.Database = h
	return nil
}

func (client *ClientCfg) Decode(v url.Values) error {
	err := sp.Decode(v, client)
	if err != nil {
		return err
	}
	client.URL, err = parseURL(client.DB)
	if err != nil {
		return err
	}
	return nil
}

func (client *ClientCfg) Encode() (url.Values, error) {
	v, err := sp.Encode(client)
	if err != nil {
		return nil, err
	}
	return v, nil
}
