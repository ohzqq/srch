package idx

import (
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"slices"

	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

type Idx struct {
	*db.DB
	Params *param.Params
}

func New(params string, data InitDB) (*Idx, error) {
	idx := Init(params)

	db, err := data()
	if err != nil {
		return idx, err
	}
	idx.DB = db

	return idx, nil
}

func NewIdx() *Idx {
	db, _ := db.New()
	return &Idx{
		Params: param.New(),
		DB:     db,
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

	var err error
	if idx.Params.Has(param.Path) {
		if ext := filepath.Ext(idx.Params.Path); ext != "" {
			idx.DB, err = NewRam(idx.Params)
			if err != nil {
				return nil, err
			}
		} else {
			idx.DB, err = OpenDisk(idx.Params)
			if err != nil {
				return nil, err
			}
		}
	}
	return idx, nil
}

func (idx *Idx) createTable(params *param.Params) error {
	err := idx.DB.CreateTable(params.IndexName)
	if err != nil {
		return err
	}
	m := idx.getDocMapping(params)
	_, err = idx.DB.CfgTable(params.IndexName, m)
	if err != nil {
		return err
	}
	return nil
}

func (idx *Idx) storageType() string {
	ext := filepath.Ext(idx.Params.Path)
	if ext == "" {
		return "disk"
	}
	return "ram"
}

func (idx *Idx) getDocMapping(params *param.Params) doc.Mapping {
	if !params.Has(param.SrchAttr) && !params.Has(param.FacetAttr) {
		return doc.DefaultMapping()
	}
	return NewMappingFromParams(params)
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
	cfg := idx.DB.GetCfg(idx.Params.IndexName)
	doc := doc.New()
	for ana, attrs := range cfg.Mapping {
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

func NewMappingFromParams(params *param.Params) doc.Mapping {
	if !params.Has(param.SrchAttr) && !params.Has(param.FacetAttr) {
		return doc.DefaultMapping()
	}

	m := doc.NewMapping()

	for _, attr := range params.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range params.FacetAttr {
		m.AddKeywords(attr)
	}

	return m
}
