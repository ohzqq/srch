package srch

import (
	"errors"
	"log"
	"os"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/blv"
	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/fuzz"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Lshortfile)
	viper.SetDefault("workers", 1)
	viper.SetDefault(param.HitsPerPage, 25)
}

type Indexer interface {
	Index(uid string, data map[string]any) error
	Batch(data []map[string]any) error
	Len() int
	Searcher
}

type Searcher interface {
	Search(query string) ([]map[string]any, error)
}

// Index is a structure for facets and data.
type Index struct {
	Indexer
	*data.Data

	data   []map[string]any
	res    *roaring.Bitmap
	isMem  bool
	Params *param.Params
}

var NoDataErr = errors.New("no data")

func newIndex() *Index {
	return &Index{
		Params: param.New(),
	}
}

func New(settings string) (*Index, error) {
	idx := newIndex()
	var err error
	idx.Params, err = param.Parse(settings)
	if err != nil {
		return nil, err
	}

	idx.Data = data.New(idx.Params.Route, idx.Params.Path)
	println(idx.Data.Route)

	switch idx.Data.Route {
	case param.Blv:
		idx.isMem = true
		idx.Indexer = blv.Open(idx.Params)
		return idx, nil
	case param.Dir, param.File:
		err = idx.GetData()
		if err != nil {
			return nil, err
		}
		idx.Indexer = fuzz.Open(idx.Params)
		idx.Batch(idx.data)
		return idx, nil
	}

	return idx, NoDataErr
}

func (idx *Index) Search(params string) (*Results, error) {
	var err error

	if idx.Indexer == nil {
		idx, err = New(params)
		println(idx.Len())
		if err != nil {
			if !errors.Is(err, NoDataErr) {
				return &Results{}, err
			}
			return NewResults([]map[string]any{}, &param.Params{})
		}
		return idx.Search(params)
	}

	p, err := param.Parse(params)
	if err != nil {
		return nil, err
	}

	q := p.Query
	r, err := idx.Indexer.Search(q)
	if err != nil {
		return nil, err
	}

	p = idx.Params
	p.Query = q
	res, err := NewResults(r, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (idx *Index) Batch(data []map[string]any) error {
	return idx.Indexer.Batch(data)
	//return idx.Indexer.Batch(idx.FilterDataBySrchAttr())
}

func (idx *Index) Len() int {
	return idx.Indexer.Len()
}

func (idx *Index) Has(key string) bool {
	return idx.Params.Has(key)
}

func (idx *Index) FilterDataBySrchAttr() []map[string]any {
	if len(idx.Params.SrchAttr) == 0 {
		return idx.data
	}
	if idx.Params.SrchAttr[0] == "*" {
		return idx.data
	}

	fields := idx.Params.SrchAttr

	if idx.isMem {
		fields = append(fields, idx.Params.Facets...)
	}

	return FilterDataByAttr(idx.data, fields)
}

func FilterDataByAttr(hits []map[string]any, fields []string) []map[string]any {
	if len(fields) < 1 {
		return hits
	}
	data := make([]map[string]any, len(hits))
	for i, d := range hits {
		data[i] = lo.PickByKeys(d, fields)
	}
	return data
}

func FilterDataByID(hits []map[string]any, uids []any, uid string) []map[string]any {
	ids := cast.ToStringSlice(uids)

	fn := func(hit map[string]any, idx int) bool {
		if uid == "" {
			return slices.Contains(ids, cast.ToString(idx))
		}
		for _, id := range ids {
			if hi, ok := hit[uid]; ok {
				return cast.ToString(hi) == id
			}
		}
		return false
	}

	f := lo.Filter(hits, fn)

	return f
}

func (idx *Index) GetData() error {

	var err error
	idx.data, err = idx.Data.Decode()
	if err != nil {
		return err
	}

	//files := idx.Params.GetDataFiles()
	//err = data.Get(&d, files...)
	//if err != nil {
	//return d, err
	//}
	return nil
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
