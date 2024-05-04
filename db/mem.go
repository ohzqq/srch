package db

import (
	"slices"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/datastores/ram"
	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/srch/doc"
)

type Mem struct {
	*hare.Database
	docs []*doc.Doc
}

func NewMem() (*Mem, error) {
	r := &ram.Ram{
		Store: store.New(),
	}
	db, err := hare.New(r)
	if err != nil {
		return nil, err
	}
	return &Mem{Database: db}, nil
}

func (m *Mem) Insert(docs ...*doc.Doc) error {
	for _, doc := range docs {
		_, err := m.Find(doc.GetID())
		if err != nil {
			m.docs = append(m.docs, doc)
			return nil
		}
		m.docs = append(m.docs, doc)
		//m.docs[i] = doc
	}
	return nil
}

func (m *Mem) Find(ids ...int) ([]*doc.Doc, error) {
	if len(ids) > 0 {
		if ids[0] == -1 {
			return m.docs, nil
		}
		var docs []*doc.Doc
		for _, doc := range m.docs {
			if slices.Contains(ids, doc.GetID()) {
				docs = append(docs, doc)
			}
		}
		return docs, nil
	}
	return m.docs, nil
}

func (m *Mem) FindAll() ([]*doc.Doc, error) {
	return m.docs, nil
}

func (m *Mem) Delete(id int) error {
	slices.DeleteFunc(m.docs, func(doc *doc.Doc) bool {
		return doc.GetID() == id
	})
	return nil
}

func (m *Mem) Index(id int) int {
	for idx, doc := range m.docs {
		if doc.GetID() == id {
			return idx
		}
	}
	return 0
}
