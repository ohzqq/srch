package idx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"slices"

	"github.com/ohzqq/hare/datastores/store"
	"github.com/ohzqq/hare/dberr"
	"github.com/ohzqq/srch/client"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Idx struct {
	*client.Client
	store  *store.Store
	Params *param.Params
}

func New(params string, data InitDB) (*Idx, error) {
	idx := Init(params)

	db, err := data()
	if err != nil {
		return idx, err
	}
	idx.Client = db

	return idx, nil
}

func NewIdx() *Idx {
	db, _ := client.New()
	return &Idx{
		Params: param.New(),
		Client: db,
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
			idx.Client, err = NewRam(idx.Params)
			if err != nil {
				return nil, err
			}
		} else {
			idx.Client, err = OpenDisk(idx.Params)
			if err != nil {
				return nil, err
			}
		}
	}
	return idx, nil
}

func (idx *Idx) createTable(params *param.Params) error {
	err := idx.Client.CreateTable(params.IndexName)
	if err != nil {
		return err
	}
	m := idx.getDocMapping(params)
	err = idx.Client.CfgTable(params.IndexName, m, params.UID)
	if err != nil {
		return err
	}
	return nil
}

func (client *Idx) storageType() string {
	ext := filepath.Ext(client.Params.Path)
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
	tbl, err := idx.Client.GetTable(idx.Params.IndexName)
	if err != nil {
		return err
	}
	doc := doc.New().WithCustomID(idx.Params.UID)
	for ana, attrs := range tbl.Mapping {
		for field, val := range data {
			if id, ok := data[doc.CustomID]; ok {
				doc.ID = cast.ToInt(id)
			}
			//if ana == analyzer.Simple && slices.Equal(attrs, []string{"*"}) {
			//  doc.AddField(ana, field, val)
			//}
			if slices.Contains(attrs, field) {
				doc.AddField(ana, field, val)
			}
		}
	}
	err = tbl.Update(doc)
	if err != nil && !errors.Is(err, dberr.ErrNoRecord) {
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
