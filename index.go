package srch

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/fuzz"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Lshortfile)
	for _, key := range param.SettingParams {
		switch key {
		case param.SrchAttr:
			viper.SetDefault(key.Snake(), []string{"title"})
		case param.FacetAttr:
			viper.SetDefault(key.Snake(), []string{"tags"})
		case param.SortAttr:
			viper.SetDefault(key.Snake(), []string{"title:desc"})
		case param.UID:
			viper.SetDefault(key.Snake(), "id")
		}
	}

	for _, key := range param.SearchParams {
		switch key {
		case param.SortFacetsBy:
			viper.SetDefault(key.Snake(), "tags:count:desc")
		case param.Facets:
			viper.SetDefault(key.Snake(), []string{"tags"})
		case param.RtrvAttr:
			viper.SetDefault(key.Snake(), "*")
		case param.Page:
			viper.SetDefault(key.Snake(), 0)
		case param.HitsPerPage:
			viper.SetDefault(key.Snake(), -1)
		case param.SortBy:
			viper.SetDefault(key.Snake(), "title")
		case param.Order:
			viper.SetDefault(key.Snake(), "desc")
		}
	}

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

	Docs   []map[string]any
	res    *roaring.Bitmap
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
		return nil, fmt.Errorf("new index param parsing err: %w\n", err)
	}

	idx.Data = data.New(idx.Params.Route, idx.Params.Path)

	idx.Indexer = fuzz.Open(idx.Params)

	switch idx.Data.Route {
	case param.Blv.String():
		idx.Params.SrchAttr = []string{"*"}
		//idx.Indexer, err = blv.Open(idx.Params)
		//if err != nil {
		//return nil, fmt.Errorf("new index open bleve err: %w\n", err)
		//}
		return idx, nil
	case param.Dir.String(), param.File.String():
		err = idx.GetData()
		if err != nil {
			return nil, err
		}
		idx.Batch(idx.Docs)
		return idx, nil
	}

	return idx, nil
}

func Mem(settings string, data []map[string]any) (*Index, error) {
	idx := newIndex()
	var err error
	idx.Params, err = param.Parse(settings)
	if err != nil {
		return nil, fmt.Errorf("new index param parsing err: %w\n", err)
	}

	idx.Indexer = fuzz.Open(idx.Params)
	idx.Docs = data
	err = idx.Batch(idx.Docs)
	if err != nil {
		return nil, fmt.Errorf("doc indexing err: %w\n", err)
	}

	return idx, nil
}

func (idx *Index) Search(params string) (*Response, error) {
	var err error

	if idx.Indexer == nil {
		idx, err = New(params)
		if err != nil {
			if !errors.Is(err, NoDataErr) {
				return &Response{}, fmt.Errorf("search err: %w\n", err)
			}
			return NewResponse([]map[string]any{}, &param.Params{})
		}
		return idx.Search(params)
	}

	p, err := param.Parse(params)
	if err != nil {
		return nil, fmt.Errorf("search failed to parse %s: err %w\n", params, err)
	}

	q := p.Query
	r, err := idx.Indexer.Search(q)
	if err != nil {
		return nil, fmt.Errorf("search '%s' failed: %w\n", q, err)
	}

	res, err := NewResponse(r, p)
	if err != nil {
		return nil, fmt.Errorf("response failed with err: %w", err)
	}
	return res, nil
}

func (idx *Index) FilterDataBySrchAttr() []map[string]any {
	if len(idx.Params.SrchAttr) == 0 {
		return idx.Docs
	}
	if idx.Params.SrchAttr[0] == "*" {
		return idx.Docs
	}

	fields := idx.Params.SrchAttr

	return FilterDataByAttr(idx.Docs, fields)
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
		var has bool
		for _, id := range ids {
			if hi, ok := hit[uid]; ok {
				if cast.ToString(hi) == id {
					has = true
				}
			}
		}
		return has
	}

	f := lo.Filter(hits, fn)

	return f
}

func (idx *Index) GetData() error {
	var err error
	idx.Docs, err = idx.Data.Decode()
	if err != nil {
		return err
	}
	return nil
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
