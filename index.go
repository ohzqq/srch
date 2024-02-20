package srch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Lshortfile)
	viper.SetDefault("workers", 1)
	viper.SetDefault(HitsPerPage, 25)
}

// Index is a structure for facets and data.
type Index struct {
	fields map[string]*Field
	Data   []map[string]any
	res    *roaring.Bitmap

	*Params `json:"params"`
}

var NoDataErr = errors.New("no data")

type SearchFunc func(string) []map[string]any

type Opt func(*Index)

func New(settings any) (*Index, error) {
	idx := newIndex()
	idx.Params = ParseParams(settings)
	idx.fields = idx.Params.Fields()

	err := idx.GetData()
	if err != nil && !errors.Is(err, NoDataErr) {
		return nil, fmt.Errorf("data parsing error: %w\n", err)
	}

	return idx, nil
}

func newIndex() *Index {
	return &Index{
		fields: make(map[string]*Field),
	}
}

func (idx *Index) Index(src []map[string]any) *Index {
	idx.Data = src

	if idx.Has(SortBy) {
		idx.Sort()
	}

	for id, d := range idx.Data {
		for _, attr := range idx.SrchAttr() {
			if val, ok := d[attr]; ok {
				idx.fields[attr].Add(val, []int{id})
			}
		}
	}

	return idx
}

func (idx *Index) Get(params string) *Response {
	return idx.Search(params)
}

func (idx *Index) Post(params any) *Response {
	p := ParseSearchParamsJSON(params)
	return idx.Search(p)
}

func (idx *Index) Search(params string) *Response {
	idx.res = idx.Bitmap()
	idx.SetSearch(params)

	query := idx.Query()
	if query != "" {
		switch idx.GetAnalyzer() {
		case TextAnalyzer:
			idx.res.And(idx.FullText(query))
		case KeywordAnalyzer:
			idx.res.And(idx.FuzzySearch(query))
		}
	}

	res := idx.Response()

	if !idx.Params.HasFilters() {
		return res
	}

	filters := idx.Params.Get(FacetFilters)
	return res.Filter(filters)
}

func (idx *Index) Sort() {
	sort := idx.Params.Get(SortBy)
	var sortType string
	for _, sb := range idx.SortAttr() {
		if t, found := strings.CutPrefix(sb, sort+":"); found {
			sortType = t
		}
	}
	switch sortType {
	case "text":
		sortDataByTextField(idx.Data, sort)
	case "num":
		sortDataByNumField(idx.Data, sort)
	}
	if idx.Params.Has(Order) {
		if idx.Params.Get(Order) == "desc" {
			slices.Reverse(idx.Data)
		}
	}
}

func (idx *Index) Response() *Response {
	return NewResponse(idx.GetResults(), idx.GetParams())
}

func (idx *Index) GetParams() url.Values {
	return idx.Values()
}

func (idx *Index) FullText(q string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, field := range idx.SearchableFields() {
		bits = append(bits, field.Filter(q))
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (idx *Index) FuzzySearch(q string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, field := range idx.SearchableFields() {
		bits = append(bits, field.Fuzzy(q))
	}
	res := roaring.ParAnd(viper.GetInt("workers"), bits...)
	return res
}

func (idx Index) Bitmap() *roaring.Bitmap {
	bits := roaring.New()
	bits.AddRange(0, uint64(len(idx.Data)))
	return bits
}

func (idx Index) HasResults() bool {
	if idx.res == nil {
		return false
	}
	if idx.res.IsEmpty() {
		return false
	}
	return true
}

func (idx Index) GetResults() []map[string]any {
	if idx.HasResults() {
		var res []map[string]any
		idx.res.Iterate(func(x uint32) bool {
			res = append(res, idx.Data[int(x)])
			return true
		})
		return res
	}
	return []map[string]any{}
}

func (idx *Index) FilterID(ids ...int) *Response {
	if !idx.HasResults() {
		idx.res = roaring.New()
	}
	for _, id := range ids {
		idx.res.AddInt(id)
	}
	return idx.Response()
}

func (idx *Index) HasData() bool {
	return idx.Params.Has(DataFile) ||
		idx.Params.Has(DataDir)
}

func (idx *Index) GetData() error {
	if !idx.HasData() {
		return NoDataErr
	}
	var data []map[string]any
	var err error
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

	idx.Index(data)
	return nil
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

func (idx *Index) GetField(attr string) *Field {
	for _, f := range idx.fields {
		if attr == f.Attribute {
			return f
		}
	}
	return &Field{Attribute: attr}
}

func (idx *Index) SearchableFields() map[string]*Field {
	return idx.fields
}

func (idx *Index) UnmarshalJSON(d []byte) error {
	un := make(map[string]json.RawMessage)
	err := json.Unmarshal(d, &un)
	if err != nil {
		return err
	}

	if msg, ok := un[Query]; ok {
		var q string
		err := json.Unmarshal(msg, &q)
		if err != nil {
			return err
		}
		idx.Params.Settings = ParseQuery(q)
	}

	if msg, ok := un[Hits]; ok {
		var data []map[string]any
		err := json.Unmarshal(msg, &data)
		if err != nil {
			return err
		}
		idx.Index(data)
	}

	return nil
}

// String satisfies the fuzzy.Source interface.
func (idx *Index) String(i int) string {
	attr := idx.SrchAttr()
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
func (idx *Index) Len() int {
	return len(idx.Data)
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
