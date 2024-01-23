package srch

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
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
	Fields   []*Field
	Data     []map[string]any
	bits     *roaring.Bitmap
	res      *roaring.Bitmap
	Settings *Settings

	*Query `json:"params"`
}

type SearchFunc func(string) []map[string]any

type Opt func(*Index)

func New(settings any) *Index {
	idx := &Index{
		Query: NewQuery(settings),
		bits:  roaring.New(),
		res:   roaring.New(),
	}
	idx.Settings = idx.GetSettings()
	idx.Fields = idx.GetSettings().Fields()

	if idx.Query.HasData() {
		d, err := idx.Query.GetData()
		if err != nil {
			return idx
		}
		return idx.Index(d)
	}

	return idx
}

func NewIndex(query any, opts ...Opt) *Index {
	idx := &Index{
		Query: NewQuery(query),
	}

	for _, opt := range opts {
		opt(idx)
	}

	return idx
}

func (idx *Index) Index(src []map[string]any) *Index {
	idx.Data = src
	idx.bits.AddRange(0, uint64(len(src)))
	idx.res.AddRange(0, uint64(len(src)))

	if idx.Query.Params.Has("sort_by") {
		idx.Sort()
	}

	idx.Fields = IndexData(idx.Data, idx.Fields)

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

func WithFullText() Opt {
	return func(idx *Index) {
		idx.Query.Params.Set("full_text", "")
	}
}

func (idx *Index) FullText(q string) []map[string]any {
	return searchFullText(
		idx.Data,
		idx.TextFields(),
		idx.Query.Params.Get(ParamQuery),
	)
}

func (idx *Index) Search(params string) *Response {
	q := NewQuery(params)

	if query := q.Query(); query != "" {
		idx.res.And(idx.FuzzySearch(query))
	}

	if q.HasFilters() {
		filterFields(idx.res, idx.Fields, q.Params.Get(ParamFacetFilters))
	}

	data := idx.getDataByBitmap(idx.res)
	return NewResponse(data, idx.Query.Params)
}

func (idx Index) Bitmap() *roaring.Bitmap {
	return idx.bits.Clone()
}

func (idx Index) getDataByBitmap(bits *roaring.Bitmap) []map[string]any {
	if bits.IsEmpty() {
		return idx.Data
	}
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, idx.Data[int(x)])
		return true
	})
	return res
}

func (idx *Index) SearchIndex(q string) *Index {
	//idx.Values.Set(QueryField, q)
	var data []map[string]any
	switch idx.Settings.TextAnalyzer {
	case Text:
		data = idx.FullText(q)
	case Fuzzy:
		data = idx.FuzzyFind(q)
	}
	return idx.Copy().Index(data)
}

func (idx *Index) search(q string) []map[string]any {
	var data []map[string]any
	switch idx.Settings.TextAnalyzer {
	case Text:
		data = idx.FullText(q)
	case Fuzzy:
		data = idx.FuzzyFind(q)
	}
	return data
}

func (idx *Index) Filter(q any) *Index {
	filtered, err := filterFields(idx.Bitmap(), idx.Fields, cast.ToString(q))
	if err != nil {
		return idx
	}
	idx.res.And(filtered)
	i := New(idx.Query)
	i.Index(idx.getDataByBitmap(idx.res))
	return i
}

func (idx *Index) SetQuery(q url.Values) *Index {
	idx.Query = &Query{
		Params: q,
	}

	if idx.Query.Params.Has("full_text") {
		idx.Settings.TextAnalyzer = Text
	}

	idx.AddField(ParseFieldsFromValues(idx.Query.Params)...)

	data, err := idx.Query.GetData()
	if err == nil {
		return idx.Index(data)
	}

	return idx
}

func (idx *Index) Sort() {
	sortDataByField(idx.Data, idx.Query.Params.Get("sort_by"))
	if idx.Query.Params.Has("order") {
		if idx.Query.Params.Get("order") == "desc" {
			slices.Reverse(idx.Data)
		}
	}
}

func (idx *Index) Copy() *Index {
	return NewIndex(idx.Query.Params)
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

func (idx *Index) AddField(fields ...*Field) *Index {
	idx.Fields = append(idx.Fields, fields...)
	return idx
}

func (idx *Index) GetField(attr string) (*Field, error) {
	for _, f := range idx.Fields {
		if f.Attribute == attr {
			return f, nil
		}
	}
	return nil, errors.New("no such field")
}

func (idx *Index) HasFilters() bool {
	return len(idx.Filters()) > 0
}

func (idx *Index) Filters() url.Values {
	return lo.OmitByKeys(idx.Query.Params, ReservedKeys)
}

// HasFacets returns true if facets are configured.
func (idx *Index) HasFacets() bool {
	return len(idx.Facets()) > 0
}

func (idx *Index) Facets() []*Field {
	return FilterFacets(idx.Fields)
}

func (idx *Index) FacetLabels() []string {
	f := idx.Facets()
	facets := make([]string, len(f))
	for i, facet := range f {
		facets[i] = facet.Attribute
	}
	return facets
}

func (idx *Index) TextFields() []*Field {
	return FilterTextFields(idx.Fields)
}

func (idx *Index) SearchableFields() []string {
	return SearchableFields(idx.Fields)
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
		idx.SetQuery(ParseQuery(q))
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
	res := map[string]any{
		Hits:       idx.Data,
		"facets":   idx.Facets(),
		ParamQuery: idx.Query.Encode(),
	}
	return json.Marshal(res)
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

func (idx *Index) FuzzyFind(q string) []map[string]any {
	matches := fuzzy.FindFrom(q, idx)
	res := make([]map[string]any, matches.Len())
	for i, m := range matches {
		res[i] = idx.Data[m.Index]
	}
	return res
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
