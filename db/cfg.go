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

func (tbl *Table) Find(ids ...int) ([]*doc.Doc, error) {
	var docs []*doc.Doc
	switch len(ids) {
	case 0:
		return docs, nil
	case 1:
		if ids[0] == -1 {
			return tbl.FindAll()
		}
		fallthrough
	default:
		for _, id := range ids {
			doc := &doc.Doc{}
			err := tbl.Table.Find(id, doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}
}

func (tbl *Table) FindAll() ([]*doc.Doc, error) {
	ids, err := tbl.IDs()
	if err != nil {
		return nil, err
	}
	return tbl.Find(ids...)
}

func (tbl *Table) Count() int {
	ids, err := tbl.IDs()
	if err != nil {
		return 0
	}
	return len(ids)
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
