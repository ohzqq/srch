package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type Cfg struct {
	ID      int         `json:"id"`
	Table   string      `json:"table"`
	Mapping doc.Mapping `json:"mapping"`
}

func NewCfg(tbl string, m doc.Mapping) *Cfg {
	return &Cfg{
		Table:   tbl,
		Mapping: m,
		ID:      1,
	}
}

func DefaultCfg() *Cfg {
	return NewCfg("index", doc.DefaultMapping())
}

func (c *Cfg) SetID(id int) {
	c.ID = id
}

func (c *Cfg) GetID() int {
	return c.ID
}

func (c *Cfg) AfterFind(_ *hare.Database) error {
	return nil
}
