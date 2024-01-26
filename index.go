package srch

import (
	"encoding/json"
	"errors"
	"log"
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
	fields []*Field
	facets []*Field
	Data   []map[string]any
	res    *roaring.Bitmap

	*Params `json:"params"`
}

type SearchFunc func(string) []map[string]any

type Opt func(*Index)

func New(settings any) (*Index, error) {
	idx := newIndex(settings)

	if idx.Params.HasData() {
		d, err := idx.Params.GetData()
		if err != nil {
			return nil, errors.New("data parsing error")
		}
		return idx.Index(d), nil
	}

	return idx, nil
}

func newIndex(settings any) *Index {
	params := ParseParams(settings)
	return &Index{
		Params: params,
		fields: params.Fields(),
		facets: params.Facets(),
	}
}

func (idx *Index) Index(src []map[string]any) *Index {
	idx.Data = src

	if idx.Params.Values.Has("sort_by") {
		idx.Sort()
	}

	for id, d := range idx.Data {
		for i, attr := range idx.SrchAttr() {
			if val, ok := d[attr]; ok {
				idx.fields[i].Add(val, []int{id})
			}
		}
		for i, attr := range idx.FacetAttr() {
			if val, ok := d[attr]; ok {
				idx.facets[i].Add(val, []int{id})
			}
		}
	}

	return idx
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

func (idx *Index) Search(params string) *Response {
	idx.res = idx.Bitmap()
	q := ParseParams(params)

	if query := q.Query(); query != "" {
		switch idx.Params.GetAnalyzer() {
		case TextAnalyzer:
			idx.res.And(idx.FullText(query))
		case KeywordAnalyzer:
			idx.res.And(idx.FuzzySearch(query))
		}
	}

	if q.HasFilters() {
		idx.Filter(q.Values.Get(FacetFilters))
	}

	return NewResponse(idx)
}

func (idx *Index) Filter(q string) *Response {
	if idx.res == nil || idx.res.IsEmpty() {
		idx.res = idx.Bitmap()
	}
	filtered, err := filterFields(idx.res, idx.facets, q)
	if err != nil {
		return NewResponse(idx)
	}
	idx.res.And(filtered)
	return NewResponse(idx)
}

func (idx Index) Bitmap() *roaring.Bitmap {
	bits := roaring.New()
	bits.AddRange(0, uint64(len(idx.Data)))
	return bits
}

func (idx Index) GetResults() []map[string]any {
	if idx.res.IsEmpty() {
		return idx.Data
	}
	var res []map[string]any
	idx.res.Iterate(func(x uint32) bool {
		res = append(res, idx.Data[int(x)])
		return true
	})
	return res
}

func (idx *Index) Sort() {
	sortDataByField(idx.Data, idx.Params.Values.Get(SortBy))
	if idx.Params.Values.Has(Order) {
		if idx.Params.Values.Get(Order) == "desc" {
			slices.Reverse(idx.Data)
		}
	}
}

func (idx *Index) GetFacet(attr string) *Field {
	for _, f := range idx.facets {
		if attr == f.Attribute {
			return f
		}
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

func (idx *Index) Facets() []*Field {
	return idx.facets
}

func (idx *Index) FacetLabels() []string {
	return lo.Map(idx.facets, func(f *Field, _ int) string {
		return f.Attribute
	})
}

func (idx *Index) SearchableFields() []*Field {
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
		idx.Params.Values = ParseQuery(q)
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

func (idx *Index) MarshalJSON() ([]byte, error) {
	return json.Marshal(idx.StringMap())
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
