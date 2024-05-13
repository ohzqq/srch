package client

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type Cfg struct {
	*hare.Table `json:"-"`

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

func DefaultCfg() *Cfg {
	return NewCfg(defaultTbl, doc.DefaultMapping(), "")
}

func (tbl *Cfg) Find(ids ...int) ([]*doc.Doc, error) {
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
			doc.WithCustomID(tbl.CustomID)
			err := tbl.Table.Find(id, doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}
}

func (tbl *Cfg) Batch(d []byte) error {
	r := bytes.NewReader(d)
	dec := json.NewDecoder(r)
	for {
		doc := &doc.Doc{}
		if err := dec.Decode(doc); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		_, err := tbl.Table.Insert(doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tbl *Cfg) FindAll() ([]*doc.Doc, error) {
	ids, err := tbl.IDs()
	if err != nil {
		return nil, err
	}
	return tbl.Find(ids...)
}

func (tbl *Cfg) Count() int {
	ids, err := tbl.IDs()
	if err != nil {
		return 0
	}
	return len(ids)
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
	tbl, err := db.GetTable(c.Name)
	if err != nil {
		return err
	}
	c.Table = tbl
	//c.Database = db
	return nil
}
