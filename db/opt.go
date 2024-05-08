package db

import (
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/net"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/srch/doc"
)

type Opt struct {
	Func func(*DB) error
	Name string
}

func NewDisk(path string) Opt {
	fn := func(db *DB) error {
		ds, err := NewDiskStorage(path)
		if err != nil {
			return fmt.Errorf("new disk opt: %w\n", err)
		}
		return db.Init(ds)
	}
	return Opt{
		Name: "NewDisk",
		Func: fn,
	}
}

func WithDisk(path string) Opt {
	fn := func(db *DB) error {
		ds, err := OpenDisk(path)
		if err != nil {
			return err
		}
		return db.Init(ds)
	}
	return Opt{
		Name: "WithDisk",
		Func: fn,
	}
}

func WithURL(uri string, d []byte) Opt {
	fn := func(db *DB) error {
		ds, err := NewNet(uri, d)
		if err != nil {
			return err
		}
		return db.Init(ds)
	}
	return Opt{
		Name: "WithURL",
		Func: fn,
	}
}

func WithRam() Opt {
	fn := func(db *DB) error {
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
	return Opt{
		Name: "WithRam",
		Func: fn,
	}
}

func WithData(d []byte) Opt {
	fn := func(db *DB) error {
		m := map[string][]byte{
			"index": d,
		}
		ds, err := ram.NewWithTables(m)
		if err != nil {
			return err
		}
		return db.Init(ds)
	}
	return Opt{
		Name: "WithData",
		Func: fn,
	}
}

func WithDefaultCfg(tbl string, m doc.Mapping, id string) Opt {
	fn := func(db *DB) error {
		return db.CfgTable(tbl, m, id)
	}
	return Opt{
		Name: "WithDefaultCfg",
		Func: fn,
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
