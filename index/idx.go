package index

import (
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/srch/param"
)

const (
	settingsTbl = "_settings"
	defaultTbl  = "default"
)

type Idx struct {
	*hare.Database
	*param.Params

	Cfg    *Cfg
	Tables map[string]int
}

func New(opts ...Opt) (*Idx, error) {
	db := &Idx{
		Tables: make(map[string]int),
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

func (idx *Idx) initDB() error {
	if idx.Database == nil {
		err := idx.memDB()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Idx) Settings() (*Settings, error) {
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

	s := NewSettings()
	s.Table = tbl

	ids, err := s.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		cfg := &Cfg{}
		err := db.Database.Find(settingsTbl, id, cfg)
		if err != nil {
			return nil, err
		}
		s.Tables[cfg.Name] = cfg
	}

	return s, nil
}

//func (idx *Idx) GetCfg(name string) (*Cfg, error) {
//  cfg := &Cfg{}
//}

func (db *Idx) setDB(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	db.Database = h
	return nil
}

func (idx *Idx) memDB() error {
	r := &ram.Ram{
		Store: store.New(),
	}
	return idx.setDB(r)
}
