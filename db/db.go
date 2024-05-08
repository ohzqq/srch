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
		_, err := db.cfg.Insert(DefaultTable())
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
	err := db.cfg.Find(id, tbl)
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

	cfg := NewCfg(name, m, id)

	var err error
	tblID := 1
	if db.TableExists(name) {
		err = db.cfg.Update(cfg)
		if err != nil {
			tblID, err = db.cfg.Insert(cfg)
			if err != nil {
				return err
			}
		}
	} else {
		tblID, err = db.cfg.Insert(cfg)
		if err != nil {
			return err
		}
	}

	db.Tables[name] = tblID

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
