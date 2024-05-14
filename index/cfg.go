package index

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type Settings struct {
	*hare.Table
	Tables map[string]*Cfg
}

func NewSettings() *Settings {
	return &Settings{
		Tables: make(map[string]*Cfg),
	}
}

type Cfg struct {
	ID       int         `json:"_id"`
	Name     string      `json:"name"`
	CustomID string      `json:"customID,omitempty"`
	Mapping  doc.Mapping `json:"mapping"`
}

func NewCfg(tbl string, m doc.Mapping, id string) *Cfg {
	return &Cfg{
		Mapping:  m,
		CustomID: id,
		ID:       1,
		Name:     tbl,
	}
}

//func (cfg *Settings) Find(name string) (*Cfg, error) {
//  cfg := &Cfg{}
//}

func DefaultCfg() *Cfg {
	return NewCfg(defaultTbl, doc.DefaultMapping(), "")
}

func (tbl *Cfg) WithCustomID(name string) *Cfg {
	tbl.CustomID = name
	return tbl
}

func (c *Cfg) SetID(id int) {
	c.ID = id
}

func (c *Cfg) GetID() int {
	return c.ID
}

func (c *Cfg) AfterFind(db *hare.Database) error {
	//println("after find")
	return nil
}
