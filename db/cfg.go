package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type Table struct {
	*hare.Table `json:"-"`
	//*hare.Database

	ID       int         `json:"_id"`
	Name     string      `json:"name"`
	CustomID string      `json:"customID,omitempty"`
	Mapping  doc.Mapping `json:"mapping"`
}

func NewCfg(tbl string, m doc.Mapping, id string) *Table {
	return &Table{
		Mapping:  m,
		CustomID: id,
		ID:       1,
		Name:     tbl,
	}
}

func DefaultTable() *Table {
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

func (c *Table) AfterFind(db *hare.Database) error {
	//println("after find")
	tbl, err := db.GetTable(c.Name)
	if err != nil {
		return err
	}
	c.Table = tbl
	//c.Database = db
	return nil
}
