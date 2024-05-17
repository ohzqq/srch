package index

import (
	"fmt"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/hare/dberr"
	"github.com/ohzqq/srch/param"
)

type ClientCfg struct {
	*hare.Table
	*param.Cfg
}

func NewClientCfg(settings any) (*ClientCfg, error) {
	cfg := &ClientCfg{
		Cfg: param.NewCfg(),
	}
	err := param.Decode(settings, cfg)
	if err != nil {
		return nil, fmt.Errorf("param decoding error: %w\n", err)
	}
	return cfg, nil
}

func (cfg *ClientCfg) GetIdxCfg(name string) (*IdxCfg, error) {
	ids, err := cfg.IDs()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		idx := &IdxCfg{}
		err := cfg.Find(id, idx)
		if err != nil {
			return nil, err
		}
		if idx.Index == name {
			return idx, nil
		}
	}

	return nil, dberr.ErrNoTable
}

func (cfg *ClientCfg) SetTbl(tbl *hare.Table) *ClientCfg {
	cfg.Table = tbl
	return cfg
}

func (cfg *ClientCfg) Datastorage() (hare.Datastorage, error) {
	return NewDatastorage(cfg.URL)
}
