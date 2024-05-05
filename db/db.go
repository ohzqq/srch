package db

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type DB struct {
	*hare.Database

	Name   string
	Tables []string
	ids    []int
}

func Open(ds hare.Datastorage) (*DB, error) {
	db := new()

	err := db.Init(ds)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func new() *DB {
	return &DB{
		Name: "index",
	}
}

func New(opts ...Opt) (*DB, error) {
	db := new()

	for _, opt := range opts {
		err := opt(db)
		if err != nil {
			return nil, err
		}
	}

	if db.Database == nil {
		ds, err := NewMem()
		if err != nil {
			return nil, err
		}
		err = db.Init(ds)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *DB) Init(ds hare.Datastorage) error {
	h, err := hare.New(ds)
	if err != nil {
		return err
	}
	db.Tables = ds.TableNames()
	db.Database = h
	return nil
}

func (db *DB) Batch(name string, d []byte) error {
	r := bytes.NewReader(d)
	dec := json.NewDecoder(r)
	for {
		doc := &doc.Doc{}
		if err := dec.Decode(doc); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err := db.Insert(name, doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Insert(name string, doc *doc.Doc) error {
	id, err := db.Database.Insert(name, doc)
	if err != nil {
		return err
	}
	db.ids = append(db.ids, id)
	return nil
}

func (db *DB) Find(name string, ids ...int) ([]*doc.Doc, error) {
	var docs []*doc.Doc
	switch len(ids) {
	case 0:
		return docs, nil
	case 1:
		if ids[0] == -1 {
			return db.FindAll(name)
		}
		fallthrough
	default:
		for _, id := range ids {
			doc := &doc.Doc{}
			err := db.Database.Find(name, id, doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}
}

func (db *DB) FindAll(name string) ([]*doc.Doc, error) {
	ids, err := db.IDs(name)
	if err != nil {
		return nil, err
	}
	return db.Find(name, ids...)
}
