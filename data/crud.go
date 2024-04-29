package data

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/ram"
)

type DB struct {
	*hare.Database
}

func NewDB() (*DB, error) {
	ds, err := ram.New(make(map[string]map[int]string))
	//ds, err := disk.New(hareTestDB, ".json")
	if err != nil {
		return nil, err
	}
	h, err := hare.New(ds)
	if err != nil {
		return nil, err
	}
	db := &DB{
		Database: h,
	}
	err = db.CreateTable("index")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *Data) Read(id string, data any) error {
	return nil
}

func (d *DB) Insert(rec hare.Record) (int, error) {
	return d.Database.Insert("index", rec)
}

func (d *Data) Update(id string, data any) error {
	return nil
}

func (d *Data) Delete(id string, data any) error {
	return nil
}
