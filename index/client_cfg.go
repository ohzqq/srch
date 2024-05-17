package index

import (
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/dberr"
	"github.com/ohzqq/srch/param"
)

type ClientCfg struct {
	tbl *hare.Table
	*param.Cfg
}

func NewClientCfg(params *param.Cfg) *ClientCfg {
	cfg := &ClientCfg{
		Cfg: params,
	}
	return cfg
}

func (cfg *ClientCfg) Insert(idx *IdxCfg) error {
	_, err := cfg.tbl.Insert(idx)
	if err != nil {
		return fmt.Errorf("cfg.Insert error\n%w\n", err)
	}
	return nil
}

func (cfg *ClientCfg) Find(name string) (*IdxCfg, error) {
	ids, err := cfg.tbl.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		idx := &IdxCfg{}
		err := cfg.tbl.Find(id, idx)
		if err != nil {
			return nil, err
		}
		if idx.Index == name {
			return idx, nil
		}
	}

	return nil, dberr.ErrNoTable
}

func (cfg *ClientCfg) Tables() ([]*IdxCfg, error) {
	ids, err := cfg.tbl.IDs()
	if err != nil {
		return nil, err
	}

	tbls := make([]*IdxCfg, len(ids))

	for i, id := range ids {
		idx := &IdxCfg{}
		err := cfg.tbl.Find(id, idx)
		if err != nil {
			return nil, err
		}
		tbls[i] = idx
	}
	return tbls, nil
}

func (cfg *ClientCfg) SetTbl(tbl *hare.Table) *ClientCfg {
	cfg.tbl = tbl
	return cfg
}
