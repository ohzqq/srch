package client

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/dberr"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
	"golang.org/x/exp/maps"
)

const (
	settingsTbl = "_settings"
	defaultTbl  = "default"
)

type Client struct {
	*hare.Database
	*param.Params

	cfg    *Cfg
	Tables map[string]int
}

func New(opts ...Opt) (*Client, error) {
	db := &Client{
		Tables: make(map[string]int),
		Params: param.New(),
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

func Open(ds hare.Datastorage) (*Client, error) {
	db := &Client{}

	err := db.Init(ds)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (client *Client) ParseParams(settings string) error {
	params, err := param.Parse(settings)
	if err != nil {
		return err
	}
	client.Params = params
	return nil
}

func (db *Client) Init(ds hare.Datastorage) error {
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

	return nil
}

func (db *Client) setDB(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	db.Database = h
	return nil
}

func (db *Client) getCfg() error {
	//Create settings table if it doesn't exist
	if !db.Database.TableExists(settingsTbl) {
		err := db.Database.CreateTable(settingsTbl)
		if err != nil {
			return fmt.Errorf("db.getCfg CreateTable error\n%w\n", err)
		}
		_, err = db.Database.Insert(settingsTbl, DefaultCfg())
		if err != nil {
			return fmt.Errorf("db.getCfg Insert error\n%w\n", err)
		}
	}

	//Get all table names and ids
	ids, err := db.Database.IDs(settingsTbl)
	if err != nil {
		return fmt.Errorf("db.getCfg IDs error\n%w\n", err)
	}
	for _, id := range ids {
		_, err := db.findTable(id)
		if err != nil {
			return fmt.Errorf("db.getCfg findTable error\n%w\n", err)
		}
	}
	return nil
	//return db.setCfg(false)
}

func (db *Client) GetCfg(name string) (*Cfg, error) {
	if tbl, ok := db.Tables[name]; ok {
		return db.findTable(tbl)
	}
	return db.findTable(-1)
}

func (db *Client) findTable(id int) (*Cfg, error) {
	tbl := &Cfg{}
	err := db.Database.Find(settingsTbl, id, tbl)
	if err != nil {
		switch {
		case errors.Is(err, dberr.ErrNoTable):
			err := db.CreateTable(tbl.Name)
			if err != nil && !errors.Is(err, dberr.ErrTableExists) {
				return nil, fmt.Errorf("db.findTable create table error:\n%w: %v\n", err, tbl.Name)
			}
		default:
			return nil, fmt.Errorf("db.findTable error:\n%w: %v\n", err, tbl.Name)
		}
	}

	db.Tables[tbl.Name] = tbl.GetID()

	return tbl, nil
}

func (db *Client) CreateTable(name string) error {
	if !db.TableExists(name) {
		err := db.Database.CreateTable(name + "-srch")
		if err != nil {
			return err
		}
		err = db.Database.CreateTable(name + "-idx")
		if err != nil {
			return err
		}
		return nil
	}
	return dberr.ErrTableExists
}

func (client *Client) storageType() string {
	ext := filepath.Ext(client.Params.Path)
	if ext == "" {
		return "disk"
	}
	return "ram"
}

func (client *Client) getDocMapping(params *param.Params) doc.Mapping {
	if !params.Has(param.SrchAttr) && !params.Has(param.FacetAttr) {
		return doc.DefaultMapping()
	}
	return NewMappingFromParams(params)
}

func (db *Client) DropTable(name string) error {
	if db.TableExists(name) {
		err := db.Database.DropTable(name + "-srch")
		if err != nil {
			return err
		}
		err = db.Database.DropTable(name + "-idx")
		if err != nil {
			return err
		}
		delete(db.Tables, name)
		return nil
	}
	return dberr.ErrNoTable
}

func (db *Client) ListTables() []string {
	return maps.Keys(db.Tables)
}

func (db *Client) TableExists(name string) bool {
	if _, ok := db.Tables[name]; ok {
		return true
	}
	return false
}

func (db *Client) CfgTable(name string, m doc.Mapping, id string) error {

	cfg := NewCfg(name, m, id)

	var err error
	tblID := 1
	if db.Database.TableExists(name) {
		err = db.Database.Update(settingsTbl, cfg)
		if err != nil {
			tblID, err = db.Database.Insert(settingsTbl, cfg)
			if err != nil {
				return err
			}
		}
	} else {
		tblID, err = db.Database.Insert(settingsTbl, cfg)
		if err != nil {
			return err
		}
	}

	db.Tables[name] = tblID

	return nil
}

func NewMappingFromParams(params *param.Params) doc.Mapping {
	if !params.Has(param.SrchAttr) && !params.Has(param.FacetAttr) {
		return doc.DefaultMapping()
	}

	m := doc.NewMapping()

	for _, attr := range params.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range params.FacetAttr {
		m.AddKeywords(attr)
	}

	return m
}
