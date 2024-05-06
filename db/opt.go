package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/net"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/srch/doc"
)

type Opt func(*DB) error

func NewDisk(path string) Opt {
	return func(db *DB) error {
		ds, err := NewDiskStorage(path)
		if err != nil {
			return err
		}
		return db.Init(ds)
	}
}

func WithDisk(path string) Opt {
	return func(db *DB) error {
		ds, err := OpenDisk(path)
		if err != nil {
			return err
		}
		return db.Init(ds)
	}
}

func WithURL(uri string, d []byte) Opt {
	return func(db *DB) error {
		ds, err := NewNet(uri, d)
		if err != nil {
			return err
		}
		return db.Init(ds)
	}
}

func WithRam(db *DB) error {
	ds, err := NewMem()
	if err != nil {
		return err
	}
	err = db.Init(ds)
	if err != nil {
		return err
	}
	return nil
}

func WithData(d []byte) Opt {
	return func(db *DB) error {
		m := map[string][]byte{
			"index": d,
		}
		ds, err := ram.New(m)
		if err != nil {
			return err
		}
		return db.Init(ds)
	}
}

func WithDefaultCfg(tbl string, m doc.Mapping) Opt {
	return func(db *DB) error {
		if db.TableExists(settingsTbl) {
			cfg := NewCfg(tbl, m)
			err := db.Update(settingsTbl, cfg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func NewDiskStore(path string) (hare.Datastorage, error) {
	return NewDiskStorage(path)
}

func NewRamStore(path string) (hare.Datastorage, error) {
	return NewMem()
}

func NewNetStore(path string, d []byte) (hare.Datastorage, error) {
	return NewNet(path, d)
}

func OpenDisk(path string) (*disk.Disk, error) {
	ds, err := disk.New(path, ".json")
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func NewDiskStorage(path string) (*disk.Disk, error) {
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
	if !ds.TableExists("index-settings") {
		err = ds.CreateTable("index-settings")
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
