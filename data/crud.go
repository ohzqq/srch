package data

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/srch/param"
)

type DB struct {
	*hare.Database
	onDisk bool
	docs   []*Doc
	Name   string
	uid    string
	*param.Params
}

func NewDB(params string, opts ...Opt) (*DB, error) {
	db := &DB{
		Name: "index",
		uid:  "id",
	}
	for _, opt := range opts {
		err := opt(db)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func NewDiskDB(path string) (*DB, error) {
	db, err := NewDB("", WithHare(path))
	if err != nil {
		return nil, err
	}
	err = db.CreateTable("index")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenHare(path string) (*hare.Database, error) {
	ds, err := disk.New(path, ".json")
	if err != nil {
		return nil, err
	}
	h, err := hare.New(ds)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (d *Data) Read(id string, data any) error {
	return nil
}

func (db *DB) Insert(doc *Doc) (int, error) {
	id := len(db.docs)
	db.docs = append(db.docs, doc)
	if db.onDisk {
		return db.Database.Insert(db.Name, doc)
	}
	return id, nil
}

func (db *DB) Find(id int) (*Doc, error) {
	doc := &Doc{}
	err := db.Database.Find(db.Name, id, doc)
	return doc, err
}

func (db *DB) Search(kw string) ([]int, error) {
	var ids []int

	docs, err := db.AllRecords()
	if err != nil {
		return ids, err
	}

	for _, doc := range docs {
		if doc.SearchAllFields(kw) {
			ids = append(ids, doc.ID)
		}
	}

	return ids, nil
}

func (db *DB) AllRecords() ([]*Doc, error) {
	ids, err := db.IDs(db.Name)
	if err != nil {
		return nil, err
	}
	docs := make([]*Doc, len(ids))
	for i, id := range ids {
		doc, err := db.Find(id)
		if err != nil {
			return nil, err
		}
		docs[i] = doc
	}
	return docs, nil
}

func (d *Data) Update(id string, data any) error {
	return nil
}

func (d *Data) Delete(id string, data any) error {
	return nil
}
