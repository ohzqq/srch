package srch

import (
	"errors"
	"log"
	"os"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/blv"
	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/mem"
	"github.com/ohzqq/srch/param"
	"github.com/ohzqq/srch/txt"
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
	Searcher
}

type Searcher interface {
	Search(query string) ([]map[string]any, error)
}

// Index is a structure for facets and data.
type Index struct {
	fields map[string]*txt.Field
	Data   []map[string]any
	res    *roaring.Bitmap
	//idx     *FullText
	isBleve bool
	idx     Indexer

	Params *param.Params
}

var NoDataErr = errors.New("no data")

type SearchFunc func(string) []map[string]any

type Opt func(*Index) error

func newIndex() *Index {
	return &Index{
		fields: make(map[string]*txt.Field),
		Params: param.New(),
		idx:    mem.New(),
	}
}

func New(settings string) (*Index, error) {
	idx := newIndex()
	var err error
	idx.Params, err = param.Parse(settings)
	if err != nil {
		return nil, err
	}

	for _, attr := range idx.Params.SrchAttr {
		idx.fields[attr] = txt.NewField(attr)
	}

	if idx.Params.SrchCfg.BlvPath != "" {
		idx.isBleve = true
		idx.idx = blv.Open(idx.Params.SrchCfg)
		return idx, nil
	}

	if idx.Params.HasData() {
		idx.Data, err = idx.GetData()
		if err != nil {
			return nil, NoDataErr
		}
		idx.idx.Batch(idx.Data)
	}

	return idx, nil
}

func (idx *Index) Search(query string) ([]map[string]any, error) {
	return idx.idx.Search(query)
}

//func (idx Idx) Bitmap() *roaring.Bitmap {
//  bits := roaring.New()

//  if idx.isBleve {
//    b, err := idx.idx.Search("")
//    if b != nil {
//      return bits
//    }
//    return b
//  }

//  if uid := idx.Params.SrchCfg.UID; uid != "" {
//    for _, d := range idx.Data {
//      bits.AddInt(cast.ToInt(d[uid]))
//    }
//  } else {
//    bits.AddRange(0, uint64(len(idx.Data)))
//  }
//  return bits
//}

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}

func (idx *Index) GetData() ([]map[string]any, error) {
	var d []map[string]any
	var err error

	files := idx.Params.GetDataFiles()
	err = data.Get(&d, files...)
	if err != nil {
		return d, err
	}
	return d, nil
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
