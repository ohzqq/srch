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
	*param.Params
}

func New(opts ...Opt) (*Client, error) {
	db := &Client{
		Params: param.New(),
	}

	for _, opt := range opts {
		err := opt(db)
		if err != nil {
			return nil, fmt.Errorf("option %v error: %w\n", opt, err)
		}
	}

	err := db.initDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (idx *Client) initDB() error {
	if idx.Database == nil {
		err := idx.memDB()
		if err != nil {
			return err
		}
	}

	return nil
}

func (idx *Client) GetCfg(name string) (*Cfg, error) {
	s, err := idx.Cfg()
	if err != nil {
		return nil, err
	}
	if cfg, ok := s[name]; ok {
		return cfg, nil
	}
	return nil, dberr.ErrNoTable
}

func (db *Client) Cfg() (map[string]*Cfg, error) {
	//Create settings table if it doesn't exist
	if !db.Database.TableExists(settingsTbl) {
		err := db.Database.CreateTable(settingsTbl)
		if err != nil {
			return nil, fmt.Errorf("db.getCfg CreateTable error\n%w\n", err)
		}
		_, err = db.Database.Insert(settingsTbl, DefaultCfg())
		if err != nil {
			return nil, fmt.Errorf("db.getCfg Insert error\n%w\n", err)
		}
	}

	tbl, err := db.Database.GetTable(settingsTbl)
	if err != nil {
		return nil, err
	}

	ids, err := tbl.IDs()
	if err != nil {
		return nil, err
	}

	tbls := make(map[string]*Cfg)
	for _, id := range ids {
		cfg := &Cfg{}
		err := db.Database.Find(settingsTbl, id, cfg)
		if err != nil {
			return nil, err
		}
		tbls[cfg.Name] = cfg
	}

	return tbls, nil
}

//func (idx *Idx) GetCfg(name string) (*Cfg, error) {
//  cfg := &Cfg{}
//}

func (db *Client) setDB(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	db.Database = h
	return nil
}

func (idx *Client) memDB() error {
	r := &ram.Ram{
		Store: store.New(),
	}
	return idx.setDB(r)
}
