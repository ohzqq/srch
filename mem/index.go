package mem

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ohzqq/srch/param"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cast"
)

type Index struct {
	Data []map[string]any
	data []string

	*param.SrchCfg
}

var NoDataErr = errors.New("no data")

func New(cfg *param.SrchCfg) *Index {
	return &Index{
		SrchCfg: cfg,
	}
}

func Open(cfg *param.SrchCfg) (*Index, error) {
	idx := New(cfg)

	data, err := idx.GetData()
	if err != nil {
		return idx, fmt.Errorf("data parsing error: %w\n", err)
	}

	err = idx.Batch(data)
	if err != nil {
		return idx, err
	}

	return idx, nil
}

func (idx *Index) Search(query string) ([]map[string]any, error) {
	if query == "" {
		return idx.Data, nil
	}
	matches := fuzzy.FindNoSort(query, idx.data)
	res := make([]map[string]any, matches.Len())
	for i, m := range matches {
		res[i] = idx.Data[m.Index]
	}
	return res, nil
}

func (idx *Index) Index(_ string, data map[string]any) error {
	idx.Data = append(idx.Data, data)

	var val []string
	for _, f := range idx.SrchAttr {
		if f == "*" {
			for _, v := range data {
				val = append(val, cast.ToString(v))
			}
		} else {
			if v, ok := data[f]; ok {
				val = append(val, cast.ToString(v))
			}
		}
	}
	idx.data = append(idx.data, strings.Join(val, " "))
	return nil
}

func (idx *Index) Batch(data []map[string]any) error {

	if len(data) < 1 {
		return NoDataErr
	}

	for _, d := range data {
		idx.Index("", d)
	}
	return nil
}

func getSrchValue(data map[string]any, fields ...string) string {

	var val []string
	for _, f := range fields {
		if v, ok := data[f]; ok {
			val = append(val, cast.ToString(v))
		}
	}
	return strings.Join(val, " ")
}

func (idx *Index) GetData() ([]map[string]any, error) {
	var data []map[string]any

	if !idx.HasData() {
		return data, NoDataErr
	}

	var err error

	files := idx.GetDataFiles()
	err = GetData(&data, files...)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (idx *Index) Len() int {
	return len(idx.data)
}

func (idx *Index) String(i int) string {
	if i < idx.Len() {
		return idx.data[i]
	}
	return ""
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
