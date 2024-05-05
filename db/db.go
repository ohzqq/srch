package db

import (
	"bytes"
	"encoding/json"
	"io"
	"slices"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/doc"
)

type DB struct {
	*hare.Database

	Name   string
	Tables []string
	ids    []int
}

func Open(ds hare.Datastorage) (*DB, error) {
	db := new()

	err := db.Init(ds)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func new() *DB {
	return &DB{
		Name: "index",
	}
}

func New(opts ...Opt) (*DB, error) {
	db := new()

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
	db.Tables = ds.TableNames()
	db.Database = h
	return nil
}

func (db *DB) Batch(d []byte) error {
	r := bytes.NewReader(d)
	dec := json.NewDecoder(r)
	for {
		doc := &doc.Doc{}
		if err := dec.Decode(doc); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err := db.Insert(doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Insert(doc *doc.Doc) error {
	id, err := db.Database.Insert(db.Name, doc)
	if err != nil {
		return err
	}
	db.ids = append(db.ids, id)
	return nil
}

//func (idx *DB) Insert(d []byte) error {
//  doc := make(map[string]any)
//  err := json.Unmarshal(d, &doc)
//  if err != nil {
//    return err
//  }
//  return idx.InsertDoc(doc)
//}

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

func (db *DB) InsertDoc(m *doc.Mapping, data map[string]any) error {
	doc := doc.New()
	for ana, attrs := range m.Mapping {
		for field, val := range data {
			if ana == analyzer.Simple && slices.Equal(attrs, []string{"*"}) {
				doc.AddField(ana, field, val)
			}
			doc.AddField(ana, field, val)
		}
	}
	id, err := db.Database.Insert(db.Name, doc)
	if err != nil {
		return err
	}
	db.ids = append(db.ids, id)
	return nil
}

//func (db *DB) Search(kw string) ([]int, error) {
//  docs, err := db.Find(-1)
//  if err != nil {
//    return nil, err
//  }

//  res := roaring.New()
//  for ana, attrs := range db.Mapping.Mapping {
//    for _, doc := range docs {
//      for _, attr := range attrs {
//        id := doc.Search(attr, ana, kw)
//        if id != -1 {
//          res.AddInt(id)
//        }
//      }
//    }
//  }
//  ids := cast.ToIntSlice(res.ToArray())
//  return ids, nil
//}
