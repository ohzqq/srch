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
	Src

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

	if db.Src == nil {
		return nil, errors.New("need a data source")
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

	err := db.Src.Insert(doc)
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

func (db *DB) Find(ids ...int) ([]*doc.Doc, error) {
	return db.Src.Find(ids...)
}

func (db *DB) Search(kw string) ([]int, error) {
	var ids []int

	docs, err := db.Find(-1)
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

func (d *Data) Update(id string, data any) error {
	return nil
}

func (d *Data) Delete(id string, data any) error {
	return nil
}
