package db

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/doc"
	"github.com/spf13/cast"
)

type DB struct {
	*hare.Database

	Name    string
	UID     string
	mapping doc.Mapping
}

func New(mapping doc.Mapping, opts ...Opt) (*DB, error) {
	db := &DB{
		Name:    "index",
		UID:     "id",
		mapping: mapping,
	}

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
	db.Database = h
	return nil
}

func (db *DB) Batch(data []map[string]any) error {
	for _, d := range data {
		err := db.Insert(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Find(ids ...int) ([]*doc.Doc, error) {
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
			err := db.Database.Find(db.Name, id, doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}
}

func (db *DB) FindAll() ([]*doc.Doc, error) {
	ids, err := db.IDs(db.Name)
	if err != nil {
		return nil, err
	}
	return db.Find(ids...)
}

func (db *DB) Insert(data map[string]any) error {
	doc := db.NewDoc(data)

	_, err := db.Database.Insert("index", doc)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) NewDoc(data map[string]any) *doc.Doc {
	return doc.New().
		SetMapping(db.mapping).
		SetData(data)
}

func (db *DB) Search(kw string) ([]int, error) {
	docs, err := db.Find(-1)
	if err != nil {
		return nil, err
	}

	res := roaring.New()
	for ana, attrs := range db.mapping {
		for _, attr := range attrs {
			if ana == analyzer.Standard {
				for _, doc := range docs {
					doc.SetMapping(db.mapping)
					id := doc.Search(attr, ana, kw)
					if id != -1 {
						res.AddInt(id)
					}
				}
			}
		}
	}
	ids := cast.ToIntSlice(res.ToArray())
	return ids, nil
}
