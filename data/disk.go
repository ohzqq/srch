package data

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/srch/doc"
)

type Disk struct {
	*hare.Database
	name string
}

func NewDisk(path string) (*Disk, error) {
	db, err := OpenHare(path)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable("index")
	if err != nil {
		return nil, err
	}
	return &Disk{
		Database: db,
	}, nil
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

func (db *Disk) Insert(doc *doc.Doc) error {
	_, err := db.Database.Insert("index", doc)
	if err != nil {
		return err
	}
	return nil
}

func (db *Disk) Find(id int) (*doc.Doc, error) {
	doc := &doc.Doc{}
	err := db.Database.Find(db.name, id, doc)
	return doc, err
}

func (db *Disk) FindAll() ([]*doc.Doc, error) {
	ids, err := db.IDs(db.name)
	if err != nil {
		return nil, err
	}
	docs := make([]*doc.Doc, len(ids))
	for i, id := range ids {
		doc, err := db.Find(id)
		if err != nil {
			return nil, err
		}
		docs[i] = doc
	}
	return docs, nil
}

func (db *Disk) Delete(id int) error {
	return db.Database.Delete(db.name, id)
}
