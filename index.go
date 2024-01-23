package srch

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
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
	fields   []*Field
	facets   []*Field
	Data     []map[string]any
	res      *roaring.Bitmap
	Settings *Settings

	*Query `json:"params"`
}

type SearchFunc func(string) []map[string]any

type Opt func(*Index)

func New(settings any) *Index {
	idx := &Index{
		Query: NewQuery(settings),
	}
	idx.Settings = idx.GetSettings()
	idx.fields = idx.GetSettings().SrchFields()
	idx.facets = idx.GetSettings().Facets()

	if idx.Query.HasData() {
		d, err := idx.Query.GetData()
		if err != nil {
			return idx
		}
		return idx.Index(d)
	}

	return idx
}

func (idx *Index) Index(src []map[string]any) *Index {
	idx.Data = src

	if idx.Query.Params.Has("sort_by") {
		idx.Sort()
	}

	if idx.GetAnalyzer() == Text {
		idx.fields = IndexData(idx.Data, idx.fields)
	}

	idx.facets = IndexData(idx.Data, idx.facets)

	return idx
}

func IndexData(data []map[string]any, fields []*Field) []*Field {
	for _, f := range fields {
		f.items = make(map[string]*FacetItem)
	}

	for id, d := range data {
		for i, f := range fields {
			if val, ok := d[f.Attribute]; ok {
				fields[i].Add(val, id)
			}
		}
	}

	return fields
}

func (idx *Index) FullText(q string) *roaring.Bitmap {
	b := FullText(idx.TextFields(), q)
	return b
}

func (idx *Index) Search(params string) *Response {
	idx.res = idx.Bitmap()
	q := NewQuery(params)
	idx.Query.Merge(q)

	if query := q.Query(); query != "" {
		switch idx.Settings.TextAnalyzer {
		case Text:
			idx.res.And(idx.FullText(query))
		case Fuzzy:
			idx.res.And(idx.FuzzySearch(query))
		}
	}

	if q.HasFilters() {
		idx.Filter(q.Params.Get(ParamFacetFilters))
	}

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

func (idx *Index) Sort() {
	sortDataByField(idx.Data, idx.Query.Params.Get("sort_by"))
	if idx.Query.Params.Has("order") {
		if idx.Query.Params.Get("order") == "desc" {
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

func (idx *Index) GetFilterValues(filters []string) map[string][]string {
	facets := make(map[string][]string)
	for _, attr := range idx.Settings.AttributesForFaceting {
		f := FilterByAttribute(attr, filters)
		if len(f) > 0 {
			facets[attr] = append(facets[attr], f...)
		}
	}
	return facets
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

func (idx *Index) TextFields() []*Field {
	return idx.fields
}

func (idx *Index) SearchableFields() []string {
	return idx.GetSrchAttr()
}

func (idx *Index) UnmarshalJSON(d []byte) error {
	un := make(map[string]json.RawMessage)
	err := json.Unmarshal(d, &un)
	if err != nil {
		return err
	}

	if msg, ok := un[ParamQuery]; ok {
		var q string
		err := json.Unmarshal(msg, &q)
		if err != nil {
			return err
		}
		idx.Query.Params = ParseQuery(q)
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
	m[ParamQuery] = idx.Query.Query()
	m[Page] = idx.Page()
	m["params"] = idx.Query
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

func (idx *Index) FuzzySearch(q string) *roaring.Bitmap {
	matches := fuzzy.FindFrom(q, idx)
	bits := roaring.New()
	for _, m := range matches {
		bits.AddInt(m.Index)
	}
	return bits
}

// String satisfies the fuzzy.Source interface.
func (idx *Index) String(i int) string {
	s := lo.PickByKeys(
		idx.Data[i],
		idx.Settings.SearchableAttributes,
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
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
