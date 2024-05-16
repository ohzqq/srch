package index

import (
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/net"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
)

var client = &http.Client{}

func NewDatastorage(u *url.URL) (hare.Datastorage, error) {
	ds, err := NewMem()
	if u == nil {
		return ds, nil
	}
	switch u.Scheme {
	case "file":
		ds, err = Disk(u.Path)
		if err != nil {
			return nil, err
		}
	case "http", "https":
		body, err := GetSettings(u)
		if err != nil {
			return nil, err
		}
		ds, err = NewNet(u.Path, body)
		if err != nil {
			return nil, err
		}
	}
	return ds, nil
}

func Disk(path string) (*disk.Disk, error) {
	p, file := filepath.Split(path)

	ext := ".json"
	if e := filepath.Ext(file); e != "" {
		ext = e
		path = p
	}
	ds, err := disk.New(path, ext)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func NewMem() (hare.Datastorage, error) {
	return DefaultMem()
}

func DefaultMem() (*ram.Ram, error) {
	r := &ram.Ram{
		Store: store.New(),
	}
	err := r.CreateTable(defaultTbl)
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

func GetSettings(u *url.URL) ([]byte, error) {
	res, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return body, nil
}
