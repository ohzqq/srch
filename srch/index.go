package srch

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/param"
	"github.com/ohzqq/srch/txt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Lshortfile)
	viper.SetDefault("workers", 1)
	viper.SetDefault(HitsPerPage, 25)
}

// Idx is a structure for facets and data.
type Idx struct {
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

type Opt func(*Idx) error

func New(settings string) (*Idx, error) {
	idx := newIndex()
	var err error
	idx.Params, err = param.Parse(settings)
	if err != nil {
		return nil, err
	}

	for _, attr := range idx.Params.SrchAttr {
		idx.fields[attr] = txt.NewField(attr)
	}

	err = idx.GetData()
	if err != nil && !errors.Is(err, NoDataErr) {
		return nil, fmt.Errorf("data parsing error: %w\n", err)
	}

	return idx, nil
}

func newIndex() *Idx {
	return &Idx{
		fields: make(map[string]*Field),
		Params: param.New(),
	}
}

func (idx Idx) Bitmap() *roaring.Bitmap {
	bits := roaring.New()
	bits.AddRange(0, uint64(len(idx.Data)))
	return bits
}

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}

func (idx *Idx) GetData() error {
	if !idx.Params.HasData() {
		return NoDataErr
	}

	var data []map[string]any
	var err error

	if idx.Params.HasData() {
		files := idx.Params.GetDataFiles()
		err = GetData(&data, files...)
		if err != nil {
			return err
		}
		idx.SetData(data)
		return nil
	}

	switch {
	case idx.Params.Has(DataFile):
		data, err = FileSrc(idx.GetSlice(DataFile)...)
		idx.Settings.Del(DataFile)
	case idx.Has(DataDir):
		data, err = DirSrc(idx.Params.Get(DataDir))
		idx.Settings.Del(DataDir)
	}
	if err != nil {
		return err
	}

	idx.SetData(data)
	return nil
}

func (idx *Idx) SetData(data []map[string]any) *Idx {
	idx.Data = data
	//return idx.Index(data)
	return idx
}

func GetDataFromQuery(q *url.Values) ([]map[string]any, error) {
	var data []map[string]any
	var err error
	switch {
	case q.Has(DataFile):
		qu := *q
		data, err = FileSrc(qu[DataFile]...)
		q.Del(DataFile)
	case q.Has(DataDir):
		data, err = DirSrc(q.Get(DataDir))
		q.Del(DataDir)
	}
	return data, err
}

// String satisfies the fuzzy.Source interface.
func (idx *Idx) String(i int) string {
	attr := idx.Params.SrchAttr
	var str string
	for _, a := range attr {
		if v, ok := idx.Data[i][a]; ok {
			str += cast.ToString(v)
			str += " "
		}
	}
	return str
}

// Len satisfies the fuzzy.Source interface.
func (idx *Idx) Len() int {
	return len(idx.Data)
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
