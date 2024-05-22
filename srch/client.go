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

	db  *hare.Database
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
	client.LoadCfg()
	return lo.Without(client.db.TableNames(), settingsTbl, "")
}

func (client *Client) TableExists(name string) bool {
	return slices.Contains(client.TableNames(), name)
}

func (client *Client) LoadCfg() error {
	if !client.db.TableExists(settingsTbl) {
		err := client.db.CreateTable(settingsTbl)
		if err != nil {
			return fmt.Errorf("client.GetCfg error\n%w\n", err)
		}
	}

	tbl, err := client.db.GetTable(settingsTbl)
	if err != nil {
		return err
	}
	client.SetTbl(tbl)

	if !client.db.TableExists(client.IndexName()) {
		err = client.db.CreateTable(client.IndexName())
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

func (client *Client) IdxIDs() ([]int, error) {
	return client.tbl.IDs()
}

func (client *Client) FindIdxCfg(name string) (*Idx, error) {
	err := client.LoadCfg()
	if err != nil {
		return nil, err
	}

	ids, err := client.tbl.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		idx := &Idx{}
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

func (client *Client) FindIdx(name string) (*Idx, error) {
	idx, err := client.FindIdxCfg(name)
	if err != nil {
		return nil, err
	}
	if client.Idx.HasSrchAttr() || client.Idx.HasFacetAttr() || client.Idx.HasSortAttr() {
		if client.HasData() {
			data := client.DataURL()
			fmt.Printf("data url %#v\n", data)
		}
	}
	if client.HasIdxURL() {
		println("2. Has idx url")
	}
	return idx, nil
}

func (client *Client) FindIdxData(name string) (*Idx, error) {
	idx, err := client.FindIdxCfg(name)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

func (client *Client) SetDatastorage(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	client.db = h
	return nil
}

func (cfg *Client) SetTbl(tbl *hare.Table) *Client {
	cfg.tbl = tbl
	return cfg
}
