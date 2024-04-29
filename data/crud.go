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
	return db, nil
}

func (d *Data) Read(id string, data any) error {
	return nil
}

func (d *Data) Create(id string, data any) error {
	return nil
}

func (d *Data) Update(id string, data any) error {
	return nil
}

func (d *Data) Delete(id string, data any) error {
	return nil
}
