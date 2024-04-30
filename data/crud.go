package data

import (
	"errors"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type DB struct {
	*hare.Database
	*param.Params

	onDisk bool
	docs   []*doc.Doc
	Name   string
}

type Src interface {
	Insert(doc ...*doc.Doc) error
	Find(id ...int) ([]*doc.Doc, error)
}

func NewDB(params string, opts ...Opt) (*DB, error) {
	p, err := param.Parse(params)
	if err != nil {
		return nil, err
	}

	db := &DB{
		Name:   "index",
		Params: p,
	}

	for _, opt := range opts {
		err := opt(db)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func (db *DB) Insert(data map[string]any) (*doc.Doc, error) {
	id := len(db.docs)
	if i, ok := data[db.UID]; ok {
		id = cast.ToInt(i)
	}

	doc := db.NewDoc(data)
	doc.SetID(id)

	err := db.insertDoc(doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (db *DB) insertDoc(doc *doc.Doc) error {
	db.docs = append(db.docs, doc)
	if db.onDisk {
		_, err := db.Database.Insert(db.Name, doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) NewDoc(data map[string]any) *doc.Doc {
	return doc.New().
		SetMapping(doc.NewMapping(db.Params)).
		SetData(data)
}

func (db *DB) Find(id int) (*doc.Doc, error) {
	if db.onDisk {
		//doc := doc.New().SetMapping(doc.NewMapping(db.Params))
		doc := &doc.Doc{}
		err := db.Database.Find(db.Name, id, doc)
		return doc, err
	}
	for _, doc := range db.docs {
		if doc.GetID() == id {
			return doc, nil
		}
	}
	return nil, errors.New("doc not found")
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

func (db *DB) AllRecords() ([]*doc.Doc, error) {
	if !db.onDisk {
		return db.docs, nil
	}
	ids, err := db.IDs(db.Name)
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

func (d *Data) Update(id string, data any) error {
	return nil
}

func (d *Data) Delete(id string, data any) error {
	return nil
}
