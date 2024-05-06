package idx

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/srch/db"
)

type Data struct {
	*db.DB
	Tables map[string]string
}

func NewData(opts ...db.Opt) (*Data, error) {
	db, err := db.New(opts...)
	if err != nil {
		return nil, err
	}
	return &Data{
		DB: db,
		Tables: map[string]string{
			"index": "",
		},
	}, nil
}

func (d *Data) CreateTable(name string) error {
}

type DataInit func() (hare.Datastorage, error)

func NewDisk(path string) (hare.Datastorage, error) {
	return NewDiskStorage(path)
}

func OpenDisk(path string) (hare.Datastorage, error) {
	return OpenDiskStorage(path)
}

func NewRam(name string) (hare.Datastorage, error) {
	return NewMemStorage(name)
}

func NewNet(name string) (hare.Datastorage, error) {
	return NewMemStorage(name)
}

func OpenDiskStorage(path string) (*disk.Disk, error) {
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

func NewMemStorage(name string) (*ram.Ram, error) {
	r := &ram.Ram{
		Store: store.New(),
	}

	err := r.CreateTable(name)
	if err != nil {
		return nil, err
	}

	err = r.CreateTable(name + "-settings")
	if err != nil {
		return nil, err
	}
	return r, nil
}

func NewNetStorage(d []byte) (*ram.Ram, error) {
	m := map[string][]byte{
		"index": d,
	}
	ds, err := ram.New(m)
	if err != nil {
		return nil, err
	}
	return ds, nil
}
