package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type Table struct {
	ID       int         `json:"id"`
	CustomID string      `json:"customID,omitempty"`
	Name     string      `json:"table"`
	Mapping  doc.Mapping `json:"mapping"`
}

func NewCfg(tbl string, m doc.Mapping, id string) *Table {
	return &Table{
		Name:     tbl,
		Mapping:  m,
		ID:       1,
		CustomID: "",
	}
}

func DefaultCfg() *Table {
	return NewCfg("index", doc.DefaultMapping(), "")
}

func (tbl *Table) WithCustomID(name string) *Table {
	tbl.CustomID = name
	return tbl
}

func (c *Table) SetID(id int) {
	c.ID = id
}

func (c *Table) GetID() int {
	return c.ID
}

func (c *Table) AfterFind(_ *hare.Database) error {
	return nil
}
