package index

import "github.com/ohzqq/hare"

type Idx struct {
	*hare.Database
	*IdxCfg
}

func NewIdx(db *hare.Database, cfg *IdxCfg) *Idx {
	return &Idx{
		Database: db,
		IdxCfg:   cfg,
	}
}
