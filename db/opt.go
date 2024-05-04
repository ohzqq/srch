package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/net"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
)

type Opt func(*DB) error

func WithUID(uid string) Opt {
	return func(db *DB) error {
		db.UID = uid
		return nil
	}
}

func WithName(name string) Opt {
	return func(db *DB) error {
		db.Name = name
		return nil
	}
}

func WithDisk(path string) Opt {
	return func(db *DB) error {
		ds, err := NewDisk(path)
		if err != nil {
			return err
		}
		h, err := hare.New(ds)
		if err != nil {
			return err
		}
		db.Database = h
		return nil
	}
}

func WithURL(uri string, d []byte) Opt {
	return func(db *DB) error {
		ds, err := NewNet(uri, d)
		if err != nil {
			return err
		}
		h, err := hare.New(ds)
		if err != nil {
			return err
		}
		db.Database = h
		return nil
	}
}

func InitDisk(path string) InitData {
	return func() (hare.Datastorage, error) {
		return NewDisk(path)
	}
}

func InitNet(uri string, d []byte) InitData {
	return func() (hare.Datastorage, error) {
		return NewNet(uri, d)
	}
}

func NewDisk(path string) (*disk.Disk, error) {
	ds, err := disk.New(path, ".json")
	if err != nil {
		return nil, err
	}
	if !ds.TableExists("index") {
		err = ds.CreateTable("index")
		if err != nil {
			return nil, err
		}
	}
	return ds, nil
}

func InitMem() (hare.Datastorage, error) {
	return NewMem()
}

func NewMem() (*ram.Ram, error) {
	r := &ram.Ram{
		Store: store.New(),
	}
	err := r.CreateTable("index")
	if err != nil {
		return nil, err
	}
	return r, nil
}

func NewNet(uri string, d []byte) (*net.Net, error) {
	ds, err := net.New(uri, d)
	if err != nil {
		return nil, err
	}
	return ds, nil
}
