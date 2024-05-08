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
	Tables map[string]*Table
}

func New(opts ...Opt) (*DB, error) {
	db := &DB{
		Tables: make(map[string]*Table),
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
	err := db.setDB(ds)
	if err != nil {
		return fmt.Errorf("db init error: %w\n", err)
	}

	err = db.getCfg()
	if err != nil {
		return fmt.Errorf("db get settings error: %w\n", err)
	}

	err = db.getTables()
	if err != nil {
		return fmt.Errorf("get tables init: %w\n", err)
	}

	//tables := db.ListTables()
	//fmt.Printf("tables list %v\n", tables)
	//switch {
	//case len(tables) == 1 && slices.Contains(tables, settingsTbl):
	//  fallthrough
	//case len(tables) == 0:
	//  err := ds.CreateTable("index")
	//  if err != nil {
	//    return err
	//  }
	//}

	//for _, name := range db.TableNames() {
	//tbl, err := db.GetTable(name)
	//if err != nil {
	//return err
	//}
	//}
	//fmt.Printf("%#v\n", ids)
	//if db.TableExists(settingsTbl) {
	//  cfg := &Table{}
	//  err := db.Database.Find(settingsTbl, 1, cfg)
	//  return fmt.Errorf("tables %v\nsettings table exists %w\n", db.TableNames(), err)
	//}

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
	cfg, err := db.GetTable(settingsTbl)
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

func (db *DB) getTable(name string) (*Table, error) {
	if tbl, ok := db.Tables[name]; ok {
		return tbl, nil
	}
	return nil, fmt.Errorf("%s: %w\n", name, dberr.ErrNoTable)
}

func (db *DB) findTable(id int) error {
	tbl := &Table{}
	err := db.cfg.Find(id, tbl)
	if err != nil {
		if errors.Is(err, dberr.ErrNoTable) {
			err = db.CreateTable(tbl.Name)
			return db.findTable(id)
		}
	}
	db.Tables[tbl.Name] = tbl
	return err
}

func (db *DB) getTables() error {
	ids, err := db.cfg.IDs()
	if err != nil {
		return fmt.Errorf("getting settings table IDs error: %w\n", err)
	}
	for _, id := range ids {
		err := db.findTable(id)
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

	return db.findTable(tblID)
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
