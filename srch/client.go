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

func (client *Client) loadSettingsTbl() error {
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

	return nil
}

func (client *Client) Indexes() map[string]*Idx {
	client.loadSettingsTbl()

	ids, err := client.tbl.IDs()
	if err != nil {
		return client.indexes
	}

	for _, id := range ids {
		idx := &Idx{}
		err := client.tbl.Find(id, idx)
		if err != nil {
			return client.indexes
		}
		client.indexes[idx.Name] = idx
	}
	return client.indexes
}

func (client *Client) HasIdx(name string) bool {
	idxs := client.Indexes()
	_, ok := idxs[name]
	return ok
}

func (client *Client) GetIdx(name string) *Idx {
	idxs := client.Indexes()
	if idx, ok := idxs[name]; ok {
		return idx
	}
	return client.Idx
}

func (client *Client) FindIdx(name string) (*Idx, error) {
	if !client.HasIdx(name) {
		// since the idx doesn't exist, insert param settings
		err := client.InsertCfg(client.Idx)
		if err != nil {
			return client.Idx, err
		}
	}

	idx, err := client.Idx.setDB(client.db)
	if err != nil {
		return idx, err
	}
	idx.SetDataURL(client.DataURL()).
		SetIdxURL(client.SrchURL())

	return idx, nil
}

func (client *Client) InsertCfg(cfg *Idx) error {
	_, err := client.tbl.Insert(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) UpdateCfg(cfg *Idx) error {
	client.Idx.SetID(cfg.GetID())
	err := client.tbl.Update(client.Idx)
	if err != nil {
		return err
	}
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
	client.loadSettingsTbl()
	return lo.Without(client.db.TableNames(), settingsTbl, "")
}

func (client *Client) TableExists(name string) bool {
	return slices.Contains(client.TableNames(), name)
}
