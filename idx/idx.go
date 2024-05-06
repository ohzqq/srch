package idx

import (
	"bytes"
	"encoding/json"
	"io"
	"slices"

	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

type Idx struct {
	*db.DB
	*doc.Mapping
	Params *param.Params
}

func New(params string, data DataInit) (*Idx, error) {
	idx := Init(params)

	ds, err := data()
	if err != nil {
		return idx, err
	}

	err = idx.DB.Init(ds)
	if err != nil {
		return nil, err
	}

	return idx, nil
}

func NewIdx() *Idx {
	db, _ := db.New()
	return &Idx{
		Params:  param.New(),
		DB:      db,
		Mapping: doc.DefaultMapping(),
	}
}

func Init(settings string) *Idx {
	idx := NewIdx()
	params, err := param.Parse(settings)
	if err != nil {
		return idx
	}
	idx.Params = params
	return idx
}

func Open(settings string) (*Idx, error) {
	idx := Init(settings)

	if idx.Params.Path != "" {
		ds, err := db.NewDiskStore(idx.Params.Path)
		if err != nil {
			return nil, err
		}

		idx.DB, err = db.Open(ds)
		if err != nil {
			return nil, err
		}
	}

	if idx.DB.TableExists(idx.Params.IndexName + "-settings") {
		err := idx.DB.Database.Find(idx.Params.IndexName+"-settings", 0, idx.Mapping)
		if err != nil {
			return nil, err
		}
	}
	return idx, nil
}

func (db *Idx) Batch(d []byte) error {
	r := bytes.NewReader(d)
	dec := json.NewDecoder(r)
	for {
		m := make(map[string]any)
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err := db.Index(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (idx *Idx) Index(data map[string]any) error {
	doc := doc.New()
	for ana, attrs := range idx.Mapping.Mapping {
		for field, val := range data {
			if ana == analyzer.Simple && slices.Equal(attrs, []string{"*"}) {
				doc.AddField(ana, field, val)
			}
			doc.AddField(ana, field, val)
		}
	}
	_, err := idx.DB.Insert(idx.Params.IndexName, doc)
	if err != nil {
		return err
	}
	return nil
}

func NewMappingFromParams(params *param.Params) *doc.Mapping {
	m := doc.NewMapping()

	for _, attr := range params.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range params.Facets {
		m.AddKeywords(attr)
	}

	return m
}
