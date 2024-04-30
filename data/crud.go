package data

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/disk"
	"github.com/ohzqq/srch/analyze"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type DB struct {
	*hare.Database
	onDisk bool
	docs   []*doc.Doc
	Name   string
	uid    string
	*param.Params
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

func NewDiskDB(path string) (*DB, error) {
	db, err := NewDB("", WithHare(path))
	if err != nil {
		return nil, err
	}
	err = db.CreateTable("index")
	if err != nil {
		return nil, err
	}
	return db, nil
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

func (d *Data) Read(id string, data any) error {
	return nil
}

func (db *DB) Insert(data map[string]any) (*Doc, error) {
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

func (db *DB) insertDoc(doc *Doc) error {
	db.docs = append(db.docs, doc)
	if db.onDisk {
		_, err := db.Database.Insert(db.Name, doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) NewDoc(data map[string]any) *Doc {
	doc := newDoc()
	for _, attr := range db.Params.SrchAttr {
		if f, ok := data[attr]; ok {
			str := cast.ToString(f)
			toks := analyze.Fulltext.Tokenize(str)
			filter := bloom.NewWithEstimates(uint(len(toks)*2), 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.Fulltext[attr] = filter
		}
	}

	for _, attr := range db.Params.Facets {
		if f, ok := data[attr]; ok {
			str := cast.ToStringSlice(f)
			toks := analyze.Keywords.Tokenize(str...)
			filter := bloom.NewWithEstimates(uint(len(toks)*5), 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.Keyword[attr] = filter
		}
	}
	return doc
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
