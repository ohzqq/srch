package index

import "github.com/ohzqq/hare"

type Idx struct {
	*hare.Database
	*Cfg
}

func NewIdx(db *hare.Database, cfg *Cfg) *Idx {
	return &Idx{
		Database: db,
		Cfg:      cfg,
	}
}
