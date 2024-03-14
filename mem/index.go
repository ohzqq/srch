package mem

import (
	"errors"
	"os"

	"github.com/ohzqq/srch/param"
	"github.com/ohzqq/srch/txt"
)

type Index struct {
	fields map[string]*txt.Field
	Data   []map[string]any

	*param.SrchCfg
}

func New(cfg *param.SrchCfg) *Index {
	return &Index{
		cfg: param.SrchCfg,
	}
}

func Open(cfg *param.SrchCfg) *Index {
	return &Index{
		SrchCfg: cfg,
	}
}

func (idx *Index) GetData() error {
	if !idx.Params.HasData() {
		return NoDataErr
	}

	var data []map[string]any
	var err error

	files := idx.Params.GetDataFiles()
	err = GetData(&data, files...)
	if err != nil {
		return err
	}
	idx.SetData(data)
	return nil
}

func (idx *Index) SetData(data []map[string]any) *Idx {
	idx.Data = data
	//return idx.Index(data)
	return idx
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
