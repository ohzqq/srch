package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/net"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
)

type Opt func(*DB)

func WithUID(uid string) Opt {
	return func(db *DB) {
		db.UID = uid
	}
}

func WithName(name string) Opt {
	return func(db *DB) {
		db.Name = name
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
