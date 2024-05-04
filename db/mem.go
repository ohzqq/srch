package db

import (
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
)

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
