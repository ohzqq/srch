package data

import (
	"github.com/ohzqq/hare/datastores/store"
)

type Mem struct {
	store *store.Store
	Name  string
}

func NewMem(idx string) (*Mem, error) {
	m := &Mem{
		Name:  idx,
		store: store.New(),
	}
	return m
}
