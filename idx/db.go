package idx

import (
	"encoding/json"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/param"
)

type Data struct {
	*db.DB
}

func NewData(path string, opts ...db.Opt) (*Data, error) {
	db, err := db.New(opts...)
	if err != nil {
		return nil, err
	}
	data := &Data{
		DB: db,
	}
	return data, nil
}

func (d *Data) New(path string) error {
	return nil
}

func WithURL(uri string) db.Opt {
	return func(db *db.DB) error {
		return nil
	}
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

func NewNetStorage(uri string) (*ram.Ram, error) {
	params, err := param.Parse(uri)
	if err != nil {
		return nil, err
	}

	if !params.Has(param.IndexName) {
		params.IndexName = "index"
	}

	r := &ram.Ram{
		Store: store.New(),
	}

	err = r.CreateTable(params.IndexName)
	if err != nil {
		return nil, err
	}

	err = r.CreateTable(params.IndexName + "-settings")
	if err != nil {
		return nil, err
	}

	m := NewMappingFromParams(params)
	d, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = r.InsertRec(params.IndexName, 1, d)
	if err != nil {
		return nil, err
	}

	return r, nil
}
