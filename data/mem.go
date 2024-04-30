package data

import (
	"slices"

	"github.com/ohzqq/srch/doc"
)

type Mem struct {
	docs []*doc.Doc
}

func (m *Mem) Insert(docs ...*doc.Doc) error {
	for i, doc := range docs {
		_, err := m.Find(doc.GetID())
		if err != nil {
			m.docs = append(m.docs, doc)
			return nil
		}
		m.docs[i] = doc
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
