package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/net"
	"github.com/ohzqq/srch/doc"
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

func (db *Net) Insert(docs ...*doc.Doc) error {
	for _, doc := range docs {
		_, err := db.Database.Insert("index", doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Net) Find(ids ...int) ([]*doc.Doc, error) {
	var docs []*doc.Doc
	switch len(ids) {
	case 0:
		return docs, nil
	case 1:
		if ids[0] == -1 {
			return db.FindAll()
		}
		fallthrough
	default:
		for _, id := range ids {
			doc := &doc.Doc{}
			err := db.Database.Find(db.name, id, doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}
}

func (db *Net) FindAll() ([]*doc.Doc, error) {
	ids, err := db.IDs(db.name)
	if err != nil {
		return nil, err
	}
	return db.Find(ids...)
}

func (db *Net) Delete(id int) error {
	return db.Database.Delete(db.name, id)
}
