package db

import (
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

const (
	settingsTbl = "_settings"
)

type DB struct {
	*hare.Database
	cfg    *hare.Table
	Tables map[string]*Table
}

func New(opts ...Opt) (*DB, error) {
	db := &DB{
		Tables: make(map[string]*Table),
	}

	for _, opt := range opts {
		err := opt.Func(db)
		if err != nil {
			return nil, fmt.Errorf("option %v error: err %w\n", opt.Name, err)
		}
	}

	if db.Database == nil {
		return New(WithRam())
	}

	return db, nil
}

func Open(ds hare.Datastorage) (*DB, error) {
	db := &DB{}

	err := db.Init(ds)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Init(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return fmt.Errorf("db init error: %w\n", err)
	}
	db.Database = h

	if !db.TableExists(settingsTbl) {
		err := db.CreateTable(settingsTbl)
		if err != nil {
			return err
		}
	}

	cfg, err := db.GetTable(settingsTbl)
	if err != nil {
		return err
	}
	db.cfg = cfg

	ids, err := cfg.IDs()
	if err != nil {
		return err
	}

	for _, id := range ids {
		cfg := &Table{}
		err := db.cfg.Find(id, cfg)
		if err != nil {
			return fmt.Errorf("db init cfg error: %w\n", err)
		}
		tbl, err := db.GetTable(cfg.Name)
		if err != nil {
			return err
		}
		cfg.Table = tbl
		fmt.Printf("%#v\n", cfg.Table)
		db.Tables[cfg.Name] = cfg
	}

	return nil
}

func (db *DB) CfgTable(name string, m doc.Mapping, id string) error {

	if !db.TableExists(settingsTbl) {
		err := db.CreateTable(settingsTbl)
		if err != nil {
			return err
		}
	}

	cfg := NewCfg(name, m, id)

	if db.TableExists(name) {
		err := db.cfg.Update(cfg)
		if err != nil {
			_, err := db.cfg.Insert(cfg)
			if err != nil {
				return err
			}
			tbl, err := db.GetTable(name)
			if err != nil {
				return err
			}
			cfg.Table = tbl
			//return fmt.Errorf("update setting error: %w\n", err)
		}
	} else {
		_, err := db.cfg.Insert(cfg)
		return err
	}
	db.Tables[name] = cfg
	return nil
}

func (db *DB) GetCfg(name string) (*Table, error) {
	if tbl, ok := db.Tables[name]; ok {
		cfg := &Table{}
		err := db.cfg.Find(tbl.GetID(), cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	return DefaultTable(), nil
}

func (db *DB) Find(name string, ids ...int) ([]*doc.Doc, error) {
	var docs []*doc.Doc
	switch len(ids) {
	case 0:
		return docs, nil
	case 1:
		if ids[0] == -1 {
			return db.FindAll(name)
		}
		fallthrough
	default:
		for _, id := range ids {
			doc := &doc.Doc{}
			err := db.Database.Find(name, id, doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}
}

func (db *DB) Count(tbl string) int {
	ids, err := db.IDs(tbl)
	if err != nil {
		return 0
	}
	return len(ids)
}

func (db *DB) FindAll(name string) ([]*doc.Doc, error) {
	ids, err := db.IDs(name)
	if err != nil {
		return nil, err
	}
	return db.Find(name, ids...)
}
