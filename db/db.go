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

	Tables []*Cfg
}

func New(opts ...Opt) (*DB, error) {
	db := &DB{}

	for _, opt := range opts {
		err := opt(db)
		if err != nil {
			return nil, fmt.Errorf("option error: err %w\n", err)
		}
	}

	if db.Database == nil {
		return New(WithRam)
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
		return err
	}
	db.Database = h

	if !db.TableExists(settingsTbl) {
		_, err := db.CfgTable("index", doc.DefaultMapping())
		if err != nil {
			return err
		}
	}

	err = db.getAllCfg()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CfgTable(name string, m doc.Mapping) (*Cfg, error) {
	if !db.TableExists(settingsTbl) {
		err := db.CreateTable(settingsTbl)
		if err != nil {
			return nil, err
		}
	}

	cfg := NewCfg(name, m)

	for i, tbl := range db.Tables {
		if tbl.Table == name {
			db.Tables[i] = cfg
			err := db.Update(settingsTbl, cfg)
			return cfg, err
		}
	}

	_, err := db.Insert(settingsTbl, cfg)
	return cfg, err
}

func (db *DB) GetCfg(tbl string) *Cfg {
	for _, cfg := range db.Tables {
		if cfg.Table == tbl {
			return cfg
		}
	}
	return DefaultCfg()
}

func (db *DB) getAllCfg() error {
	tbls, err := db.IDs(settingsTbl)
	if err != nil {
		return err
	}

	db.Tables = make([]*Cfg, len(tbls))
	for i, tbl := range tbls {
		cfg := &Cfg{}
		err := db.Database.Find(settingsTbl, tbl, cfg)
		if err != nil {
			return err
		}
		db.Tables[i] = cfg
	}
	return nil
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
