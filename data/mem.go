package data

import (
	"errors"
	"slices"

	"github.com/ohzqq/srch/doc"
)

type Mem struct {
	docs []*doc.Doc
}

func (m *Mem) Insert(doc *doc.Doc) error {
	m.docs = append(m.docs, doc)
	return nil
}

func (m *Mem) Find(id int) (*doc.Doc, error) {
	for _, doc := range m.docs {
		if doc.GetID() == id {
			return doc, nil
		}
	}
	return nil, errors.New("doc not found")
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
