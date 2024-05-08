package db

import (
	"errors"
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/dberr"
	"github.com/ohzqq/srch/doc"
	"github.com/samber/lo"
)

const (
	settingsTbl = "_settings"
	defaultTbl  = "default"
)

type DB struct {
	*hare.Database
	cfg    *Table
	Tables map[string]int
}

func New(opts ...Opt) (*DB, error) {
	db := &DB{
		Tables: make(map[string]int),
	}

	for _, opt := range opts {
		err := opt.Func(db)
		if err != nil {
			return nil, fmt.Errorf("option %v error: %w\n", opt.Name, err)
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
	//step 1: initialize hare.DB
	err := db.setDB(ds)
	if err != nil {
		return fmt.Errorf("db init error: %w\n", err)
	}

	//step 2: get the settings for all indexes
	err = db.getCfg()
	if err != nil {
		return fmt.Errorf("db get settings error: %w\n", err)
	}

	//step 3: get all tables
	err = db.getTables()
	if err != nil {
		return fmt.Errorf("get tables init: %w\n", err)
	}

	return nil
}

func (db *DB) setDB(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	db.Database = h
	return nil
}

func (db *DB) getCfg() error {
	if !db.TableExists(settingsTbl) {
		err := db.CreateTable(settingsTbl)
		if err != nil {
			return err
		}
		return db.setCfg(true)
	}

	return db.setCfg(false)
}

func (db *DB) setCfg(setDefault bool) error {
	cfg, err := db.Database.GetTable(settingsTbl)
	if err != nil {
		return err
	}
	db.cfg = &Table{
		Table: cfg,
		Name:  "_settings",
	}

	if setDefault {
		_, err := db.cfg.Table.Insert(DefaultTable())
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) GetTable(name string) (*Table, error) {
	if tbl, ok := db.Tables[name]; ok {
		return db.findTable(tbl)
	}
	return db.findTable(-1)
}

func (db *DB) findTable(id int) (*Table, error) {
	tbl := &Table{}
	err := db.cfg.Table.Find(id, tbl)
	if err != nil {
		if errors.Is(err, dberr.ErrNoTable) {
			err = db.CreateTable(tbl.Name)
			return db.findTable(id)
		}
	}
	db.Tables[tbl.Name] = tbl.GetID()
	return tbl, err
}

func (db *DB) getTables() error {
	ids, err := db.cfg.IDs()
	if err != nil {
		return fmt.Errorf("getting settings table IDs error: %w\n", err)
	}
	for _, id := range ids {
		_, err := db.findTable(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) ListTables() []string {
	return lo.Without(db.TableNames(), settingsTbl, "")
}

func (db *DB) CfgTable(name string, m doc.Mapping, id string) error {

	cfg := NewTable(name, m, id)

	var err error
	tblID := 1
	if db.TableExists(name) {
		err = db.cfg.Update(cfg)
		if err != nil {
			tblID, err = db.cfg.Table.Insert(cfg)
			if err != nil {
				return err
			}
		}
	} else {
		tblID, err = db.cfg.Table.Insert(cfg)
		if err != nil {
			return err
		}
	}

	db.Tables[name] = tblID

	return nil
}
