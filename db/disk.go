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

func Open(path string) (*Disk, error) {
	db, err := OpenHare(path)
	if err != nil {
		return nil, err
	}
	return &Disk{
		Database: db,
		name:     "index",
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

func (db *Disk) Insert(docs ...*doc.Doc) error {
	for _, doc := range docs {
		_, err := db.Database.Insert("index", doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Disk) Find(ids ...int) ([]*doc.Doc, error) {
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

func (db *Disk) FindAll() ([]*doc.Doc, error) {
	ids, err := db.IDs(db.name)
	if err != nil {
		return nil, err
	}
	return db.Find(ids...)
}

func (db *Disk) Delete(id int) error {
	return db.Database.Delete(db.name, id)
}
