package db

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/doc"
	"github.com/spf13/cast"
)

type DB struct {
	Src

	Name    string
	UID     string
	mapping doc.Mapping
}

type Src interface {
	Insert(doc ...*doc.Doc) error
	Find(id ...int) ([]*doc.Doc, error)
}

func New(src Src, mapping doc.Mapping, opts ...Opt) (*DB, error) {
	db := &DB{
		Name:    "index",
		mapping: mapping,
		Src:     src,
		UID:     "id",
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
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

func (db *DB) Insert(data map[string]any) error {
	var id int
	if i, ok := data[db.UID]; ok {
		id = cast.ToInt(i)
	}

	doc := db.NewDoc(data)
	doc.SetID(id)

	err := db.Src.Insert(doc)
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
