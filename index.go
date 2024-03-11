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
	"github.com/ohzqq/srch/blv"
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
	fields map[string]*Field
	Data   []map[string]any
	res    *roaring.Bitmap
	//idx     *FullText
	isBleve bool
	idx     Indexer

	*Params `json:"params"`
}

var NoDataErr = errors.New("no data")

type SearchFunc func(string) []map[string]any

type Opt func(*Idx) error

func New(settings any) (*Idx, error) {
	idx := newIndex()
	idx.Params = ParseParams(settings)
	idx.fields = idx.Params.Fields()

	//data := idx.Params.GetData()
	//println(data)

	err := idx.GetData()
	if err != nil && !errors.Is(err, NoDataErr) {
		return nil, fmt.Errorf("data parsing error: %w\n", err)
	}

	//blv, err := bleve.Open(blevePath)
	//idx.idx, err = NewTextIndex(FTPath(blevePath))
	//if err != nil {
	//return nil, err
	//return idx, nil
	//}
	//idx.idx = &FullText{
	//Index: blv,
	//}

	return idx, nil
}

func newIndex() *Idx {
	return &Idx{
		fields: make(map[string]*Field),
	}
}

func (idx *Idx) Index(src []map[string]any) *Idx {
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

func (idx *Idx) Get(params string) *Response {
	return idx.Search(params)
}

func (idx *Idx) Post(params any) *Response {
	p := ParseSearchParamsJSON(params)
	return idx.Search(p)
}

func (idx *Idx) Search(params string) *Response {
	idx.res = idx.Bitmap()
	idx.SetSearch(params)

	query := idx.Query()
	if query != "" {
		if idx.Params.IsFullText() {
			idx.idx = blv.Open(idx.Params.GetFullText())
			bits, err := idx.idx.Search(query)
			if err != nil {
				log.Fatal(err)
			}
			idx.res.And(bits)
			return idx.Response()
		}

		idx.res.And(idx.FuzzySearch(query))
		return idx.Response()
	}

	res := idx.Response()

	if !idx.Params.HasFilters() {
		return res
	}

	//filters := idx.Params.Get(FacetFilters)
	//return res.Filter(filters)
	return res
}

func (idx *Idx) Response() *Response {
	return NewResponse(idx.GetResults(), idx.GetParams())
}

func (idx *Idx) Sort() {
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

func (idx *Idx) GetParams() url.Values {
	return idx.Values()
}

func (idx *Idx) FuzzySearch(q string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, field := range idx.SearchableFields() {
		bits = append(bits, field.Fuzzy(q))
	}
	res := roaring.ParAnd(viper.GetInt("workers"), bits...)
	return res
}

func (idx Idx) Bitmap() *roaring.Bitmap {
	bits := roaring.New()
	bits.AddRange(0, uint64(len(idx.Data)))
	return bits
}

func (idx Idx) HasResults() bool {
	if idx.res == nil {
		return false
	}
	if idx.res.IsEmpty() {
		return false
	}
	return true
}

func (idx Idx) GetResults() []map[string]any {
	if idx.HasResults() {
		return ItemsByBitmap(idx.Data, idx.res)
	}
	return []map[string]any{}
}

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}

func (idx *Idx) FilterID(ids ...int) *Response {
	if !idx.HasResults() {
		idx.res = roaring.New()
	}
	for _, id := range ids {
		idx.res.AddInt(id)
	}
	return idx.Response()
}

func (idx *Idx) GetData() error {
	if !idx.Params.HasData() {
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

	idx.SetData(data)
	return nil
}

func (idx *Idx) SetData(data []map[string]any) *Idx {
	idx.Data = data
	return idx.Index(data)
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

func (idx *Idx) GetField(attr string) *Field {
	for _, f := range idx.fields {
		if attr == f.Attribute {
			return f
		}
	}
	return &Field{Attribute: attr}
}

func (idx *Idx) SearchableFields() map[string]*Field {
	return idx.fields
}

func (idx *Idx) UnmarshalJSON(d []byte) error {
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
func (idx *Idx) String(i int) string {
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
func (idx *Idx) Len() int {
	return len(idx.Data)
}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
