package idx

import (
	"fmt"

	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/doc"
	"github.com/ohzqq/srch/param"
)

type Idx struct {
	*db.DB
	Params *param.Params
}

func New() *Idx {
	m := doc.DefaultMapping()
	db, _ := db.New(m)
	return &Idx{
		Params: param.New(),
		DB:     db,
	}
}

func Open(settings string) (*Idx, error) {
	//idx := &Idx{}
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

	return idx, nil
}

func NewMappingFromParams(params *param.Params) doc.Mapping {
	m := make(doc.Mapping)

	for _, attr := range params.SrchAttr {
		m[analyzer.Standard] = append(m[analyzer.Standard], attr)
	}

	for _, attr := range params.FacetAttr {
		m[analyzer.Keyword] = append(m[analyzer.Keyword], attr)
	}

	if m == nil {
		return doc.DefaultMapping()
	}

	return m
}
