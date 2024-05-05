package idx

import (
	"fmt"

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
	idx := New()
	var err error
	idx.Params, err = param.Parse(settings)
	if err != nil {
		return nil, fmt.Errorf("new index param parsing err: %w\n", err)
	}

	return idx, nil
}
