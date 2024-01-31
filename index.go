package srch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Lshortfile)
	viper.SetDefault("workers", 1)
}

// Index is a structure for facets and data.
type Index struct {
	fields map[string]*Field
	facets map[string]*Field
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
	idx.facets = idx.Params.Facets()

	err := idx.GetData()
	if err != nil && !errors.Is(err, NoDataErr) {
		return nil, fmt.Errorf("data parsing error: %w\n", err)
	}

	return idx, nil
}

func newIndex() *Index {
	return &Index{
		fields: make(map[string]*Field),
		facets: make(map[string]*Field),
	}
}

func (idx *Index) Index(src []map[string]any) *Index {
	idx.Data = src

	//if idx.Params.Values.Has("sort_by") {
	//  idx.Sort()
	//}

	for id, d := range idx.Data {
		for _, attr := range idx.SrchAttr() {
			if val, ok := d[attr]; ok {
				idx.fields[attr].Add(val, []int{id})
			}
		}
		for _, attr := range idx.FacetAttr() {
			if val, ok := d[attr]; ok {
				idx.facets[attr].Add(val, []int{id})
			}
		}
	}

	return idx
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

	if idx.Has(FacetFilters) {
		filters := idx.Get(FacetFilters)
		return idx.Filter(filters)
	}

	return idx.Response()
}

func (idx *Index) Filter(q string) *Response {
	if !idx.HasResults() {
		idx.res = idx.Bitmap()
	}

	idx.Set(FacetFilters, q)

	filtered, err := Filter(idx.res, idx.facets, q)
	if err != nil {
		return idx.Response()
	}

	idx.res.And(filtered)
	return idx.Response()
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
	return idx.Data
}

func (idx *Index) Sort() {
	sortDataByField(idx.Data, idx.Get(SortBy))
	if idx.Has(Order) {
		if idx.Get(Order) == "desc" {
			slices.Reverse(idx.Data)
		}
	}
}

func (idx *Index) HasData() bool {
	return idx.Has(DataFile) ||
		idx.Has(DataDir)
}

func (idx *Index) GetData() error {
	if !idx.HasData() {
		return NoDataErr
	}
	var data []map[string]any
	var err error
	switch {
	case idx.Has(DataFile):
		data, err = FileSrc(idx.GetSlice(DataFile)...)
		idx.Settings.Del(DataFile)
	case idx.Has(DataDir):
		data, err = DirSrc(idx.Get(DataDir))
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

func (idx *Index) GetFacet(attr string) *Field {
	if f, ok := idx.facets[attr]; ok {
		return f
	}
	return &Field{Attribute: attr}
}

func (idx *Index) GetField(attr string) *Field {
	for _, f := range idx.fields {
		if attr == f.Attribute {
			return f
		}
	}
	return &Field{Attribute: attr}
}

// HasFacets returns true if facets are configured.
func (idx *Index) HasFacets() bool {
	return len(idx.facets) > 0
}

func (idx *Index) Facets() map[string]*Field {
	return idx.facets
}

func (idx *Index) FacetLabels() []string {
	return lo.Keys(idx.facets)
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

func (idx *Index) StringMap() map[string]any {
	m := make(map[string]any)
	m[Query] = idx.Params.Query()
	m[Page] = idx.Page()
	m["params"] = idx.Params
	m[HitsPerPage] = idx.HitsPerPage()
	m[ParamFacets] = idx.Facets()
	return m
}

// JSON marshals an Index to json.
func (idx *Index) JSON() []byte {
	d, err := json.Marshal(idx)
	if err != nil {
		return []byte{}
	}
	return d
}

// Print writes Index json to stdout.
func (idx *Index) Print() {
	enc := json.NewEncoder(os.Stdout)
	err := enc.Encode(idx)
	if err != nil {
		log.Fatal(err)
	}
}

// PrettyPrint writes Index indented json to stdout.
func (idx *Index) PrettyPrint() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(idx)
	if err != nil {
		log.Fatal(err)
	}
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

func sortDataByField(data []map[string]any, field string) []map[string]any {
	fn := func(a, b map[string]any) int {
		x := cast.ToString(a[field])
		y := cast.ToString(b[field])
		switch {
		case x > y:
			return 1
		case x == y:
			return 0
		default:
			return -1
		}
	}
	slices.SortFunc(data, fn)
	return data
}
