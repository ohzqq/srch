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

	return client, nil
}

func (client *Client) TableNames() []string {
	return lo.Without(client.Database.TableNames(), "_settings", "")
}

func (client *Client) TableExists(name string) bool {
	return slices.Contains(client.TableNames(), name)
}

func (client *Client) GetCfg() error {
	if !client.Database.TableExists(settingsTbl) {
		err := client.Database.CreateTable(settingsTbl)
		if err != nil {
			return fmt.Errorf("client.GetCfg error\n%w\n", err)
		}
	}

	tbl, err := client.Database.GetTable(settingsTbl)
	if err != nil {
		return err
	}
	client.SetTbl(tbl)

	if !client.TableExists(client.IndexName()) {
		err = client.Database.CreateTable(client.IndexName())
		if err != nil {
			return err
		}
		_, err = client.tbl.Insert(client.Idx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) Indexes() (map[string]int, error) {
	idxs := make(map[string]int)

	err := client.GetCfg()
	if err != nil {
		return nil, err
	}

	ids, err := client.tbl.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		idx := &IdxCfg{}
		err := client.tbl.Find(id, idx)
		if err != nil {
			return nil, fmt.Errorf("%w: %v\n", err, client.IndexName())
		}
		idxs[idx.Name] = id
	}

	return idxs, err
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
	err := client.GetCfg()
	if err != nil {
		return nil, err
	}

	ids, err := client.tbl.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		idx := &IdxCfg{}
		err := client.tbl.Find(id, idx)
		if err != nil {
			return nil, fmt.Errorf("%w: %v\n", err, client.IndexName())
		}
		if idx.Name == name {
			return idx, nil
		}
	}

	return nil, err
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
