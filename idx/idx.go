package idx

import (
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/param"
)

type Idx struct {
	DB     *db.DB
	Params *param.Params
}

func New() *Idx {
	return &Idx{
		Params: param.New(),
	}
}

func Open(settings string) (*Idx, error) {
}
