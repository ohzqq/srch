package idx

import (
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

type Idx struct {
	*db.DB
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

	//idx.Mapping = NewMappingFromParams(params)

	return idx
}

func Open(settings string) (*Idx, error) {
	idx := Init(settings)

	var ds hare.Datastorage
	var err error
	if idx.Params.Has(param.Path) {
		if ext := filepath.Ext(idx.Params.Path); ext != "" {
			tbl := strings.TrimSuffix(idx.Params.Path, ext)
			if !idx.DB.TableExists(tbl) {
				err = idx.DB.CreateTable(idx.Params.IndexName)
				if err != nil {
					return nil, err
				}
			}
		} else {
			ds, err = db.NewDiskStore(idx.Params.Path)
			if err != nil {
				return nil, err
			}
			idx.DB, err = db.Open(ds)
			if err != nil {
				return nil, err
			}
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
	m := doc.NewMapping()

	for _, attr := range params.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range params.Facets {
		m.AddKeywords(attr)
	}

	return m
}
