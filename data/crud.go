package data

import (
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/doc"
	"github.com/samber/lo"
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

func NewDB(src Src, mapping doc.Mapping, opts ...Opt) (*DB, error) {
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

	//var fields []*roaring.Bitmap
	var res [][]int
	for ana, attrs := range db.mapping {
		for _, attr := range attrs {
			var ids []int
			if ana == analyzer.Fulltext {
				for _, doc := range docs {
					id := doc.Search(attr, kw)
					if id != -1 {
						if attr == "tags" {
							println(attr)
						}
						ids = append(ids, id)
					}
				}
				res = append(res, ids)
			}
		}
	}

	var ids []int
	c := 1
	for i := 0; c < len(res); i++ {
		ids = lo.Intersect(res[i], res[c])
		c++
	}
	return ids, nil
}
