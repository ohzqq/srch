package srch

import (
	"fmt"
	"slices"
	"strings"

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

func (client *Client) idxNames() []string {
	client.LoadCfg()
	tbls := lo.Without(client.db.TableNames(), settingsTbl, "")

	var names []string
	for _, tbl := range tbls {
		if strings.HasSuffix(tbl, "Idx") {
			it := strings.TrimSuffix(tbl, "Idx")
			names = append(names, it)
		}
		if strings.HasSuffix(tbl, "Data") {
			dt := strings.TrimSuffix(tbl, "Data")
			names = append(names, dt)
		}
	}

	return names
}

func (client *Client) TableNames() []string {
	client.LoadCfg()
	return lo.Without(client.db.TableNames(), settingsTbl, "")
}

func (client *Client) TableExists(name string) bool {
	return slices.Contains(client.TableNames(), name)
}

func (client *Client) HasIdx(name string) bool {
	return false
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

	if !client.db.TableExists(client.Idx.idxTblName()) {
		err = client.db.CreateTable(client.Idx.idxTblName())
		if err != nil {
			return nil, err
		}
		_, err = client.tbl.Insert(client.Idx)
		if err != nil {
			return nil, err
		}
	}

	if !client.db.TableExists(client.Idx.dataTblName()) {
		err = client.db.CreateTable(client.Idx.dataTblName())
		if err != nil {
			return nil, err
		}
	}

	cfgs, err := client.getIdxCfgs()
	if err != nil {
		return nil, err
	}

	if idx, ok := cfgs[client.IndexName()]; ok {
		return idx, nil
	}

	return nil, err
}

func (client *Client) getIdxCfgs() (map[string]*Idx, error) {
	ids, err := client.tbl.IDs()
	if err != nil {
		return nil, err
	}

	idxs := make(map[string]*Idx)
	for _, id := range ids {
		idx := &Idx{}
		err := client.tbl.Find(id, idx)
		if err != nil {
			return nil, fmt.Errorf("%w: %v\n", err, client.IndexName())
		}
		idxs[client.IndexName()] = idx
	}

	return idxs, nil
}

func (client *Client) FindIdx(name string) (*Idx, error) {
	cfg, err := client.FindIdxCfg(name)
	if err != nil {
		return nil, err
	}
	if client.Idx.HasSrchAttr() || client.Idx.HasFacetAttr() || client.Idx.HasSortAttr() {
		if client.HasData() {
			//data := client.DataURL()
			//fmt.Printf("data url %#v\n", data)
		}
	}
	if client.HasIdxURL() {
		//println("2. Has idx url")
	}
	return cfg, nil
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
