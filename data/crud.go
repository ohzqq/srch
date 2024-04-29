package data

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/srch/param"
)

type DB struct {
	*hare.Database
	Name string
	*param.Params
}

func NewDB(name string) *DB {
	if name == "" {
		name = "index"
	}
	return &DB{Name: name}
}

func NewMemDB() (*DB, error) {
	ds, err := ram.New(make(map[string]map[int]string))
	//ds, err := disk.New(hareTestDB, ".json")
	if err != nil {
		return nil, err
	}
	h, err := hare.New(ds)
	if err != nil {
		return nil, err
	}
	db := NewDB("")
	db.Database = h
	err = db.CreateTable(db.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDiskDB(path string) (*DB, error) {
	db, err := OpenDB(path)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(db.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenDB(path string) (*DB, error) {
	ds, err := disk.New(path, ".json")
	if err != nil {
		return nil, err
	}
	h, err := hare.New(ds)
	if err != nil {
		return nil, err
	}
	db := NewDB("")
	db.Database = h
	return db, nil
}

func (d *Data) Read(id string, data any) error {
	return nil
}

func (db *DB) Insert(rec hare.Record) (int, error) {
	return db.Database.Insert(db.Name, rec)
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
