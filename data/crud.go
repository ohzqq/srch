package data

import (
	"errors"

	"github.com/ohzqq/srch/analyze"
	"github.com/ohzqq/srch/doc"
	"github.com/spf13/cast"
)

type DB struct {
	Src

	Name    string
	UID     string
	mapping map[string]analyze.Analyzer
}

type Src interface {
	Insert(doc ...*doc.Doc) error
	Find(id ...int) ([]*doc.Doc, error)
}

func NewDB(src Src, mapping map[string]analyze.Analyzer, opts ...Opt) (*DB, error) {
	db := &DB{
		Name:    "index",
		mapping: mapping,
		Src:     src,
		UID:     "id",
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

func (db *DB) Batch(data []map[string]any) error {
	for _, d := range data {
		_, err := db.Insert(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Insert(data map[string]any) (*doc.Doc, error) {
	var id int
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

func (db *DB) NewDoc(data map[string]any) *doc.Doc {
	return doc.New().
		SetMapping(db.mapping).
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
