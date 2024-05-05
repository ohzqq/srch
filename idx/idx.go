package idx

import (
	"fmt"

	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

type Idx struct {
	*db.DB
	*doc.Mapping
	Params *param.Params
}

func New() *Idx {
	db, _ := db.New()

	db.Database.Insert("index-settings", doc.DefaultMapping())
	return &Idx{
		Params:  param.New(),
		DB:      db,
		Mapping: m,
	}
}

func Open(settings string) (*Idx, error) {
	idx := New()
	params, err := param.Parse(settings)
	if err != nil {
		return nil, fmt.Errorf("new index param parsing err: %w\n", err)
	}
	idx.Params = params

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
