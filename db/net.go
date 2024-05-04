package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/net"
)

type Net struct {
	*hare.Database
	name string
}

func NewNet(uri string, d []byte) (*Net, error) {
	ds, err := net.New(uri, d)
	if err != nil {
		return nil, err
	}
	db, err := hare.New(ds)
	if err != nil {
		return nil, err
	}
	return &Net{
		Database: db,
		name:     "index",
	}, nil
}
