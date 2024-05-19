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

	tbl *hare.Table

	DB    string `json:"-" mapstructure:"path" qs:"db"`
	Index string `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	UID   string `query:"uid,omitempty" json:"uid,omitempty" mapstructure:"uid" qs:"uid"`
}

func NewClient(settings any) (*Client, error) {
	client := &Client{
		Index: "default",
		URL:   &url.URL{},
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

	return client, nil
}

func (client *Client) init() error {
	//Create settings table if it doesn't exist
	if !client.Database.TableExists(settingsTbl) {
		err := client.Database.CreateTable(settingsTbl)
		if err != nil {
			return fmt.Errorf("db.getCfg CreateTable error\n%w\n", err)
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

func (client *Client) Decode(v url.Values) error {
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

func (client *Client) Encode() (url.Values, error) {
	v, err := sp.Encode(client)
	if err != nil {
		return nil, err
	}
	return v, nil
}
