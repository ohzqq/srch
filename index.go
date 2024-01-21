package srch

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"os"
	"slices"
	"strings"

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
	search   SearchFunc
	Fields   []*Field         `json:"fields"`
	Data     []map[string]any `json:"data"`
	Values   url.Values       `json:"query"`
	Settings *Settings
	*Query
}

type SearchFunc func(string) []map[string]any

type Opt func(*Index)

func New(data []map[string]any, settings *Settings) *Index {
	idx := &Index{
		Settings: settings,
		Fields:   settings.Fields(),
	}
	return idx.Index(data)
}

func NewIndex(query any, opts ...Opt) *Index {
	idx := &Index{
		Values: ParseQuery(query),
	}

	for _, opt := range opts {
		opt(idx)
	}

	idx.SetQuery(idx.Values)

	return idx
}

func (idx *Index) Index(src []map[string]any) *Index {
	if len(idx.Fields) < 1 {
		idx.AddField(NewField("title", Text))
		idx.Values.Add("field", "title")
	}

	idx.Data = src

	if idx.Values.Has("sort_by") {
		idx.Sort()
	}

	idx.Fields = IndexData(idx.Data, idx.Fields)

	if idx.HasFilters() {
		return idx.Filter(idx.Filters())
	}

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
		idx.Values.Set("full_text", "")
	}
}

func (idx *Index) FullText(q string) []map[string]any {
	return searchFullText(idx.Data, idx.TextFields(), idx.Values.Get(QueryField))
}

func (idx *Index) Search(q string) *Index {
	idx.Values.Set(QueryField, q)
	data := idx.search(q)
	return idx.Copy().Index(data)
}

func (idx *Index) Filter(q any) *Index {
	vals, err := ParseValues(q)
	if err != nil {
		return idx
	}
	idx.Data = Filter(idx.Data, idx.Facets(), vals)
	idx.Fields = IndexData(idx.Data, idx.Fields)
	return idx
}

func (idx *Index) SetQuery(q url.Values) *Index {
	idx.Values = q

	if idx.Values.Has("full_text") {
		idx.Settings.TextAnalyzer = Text
	}

	idx.AddField(ParseFieldsFromValues(idx.Values)...)

	data, err := GetDataFromQuery(&idx.Values)
	if err == nil {
		return idx.Index(data)
	}

	return idx
}

func (idx *Index) Sort() {
	sortDataByField(idx.Data, idx.Values.Get("sort_by"))
	if idx.Values.Has("order") {
		if idx.Values.Get("order") == "desc" {
			slices.Reverse(idx.Data)
		}
	}
}

func (idx *Index) Copy() *Index {
	return NewIndex(idx.Values)
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
	return lo.OmitByKeys(idx.Values, ReservedKeys)
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

	if msg, ok := un["query"]; ok {
		var q string
		err := json.Unmarshal(msg, &q)
		if err != nil {
			return err
		}
		idx.SetQuery(ParseQuery(q))
	}

	if msg, ok := un["data"]; ok {
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
		"data":   idx.Data,
		"facets": idx.Facets(),
		"query":  idx.Values.Encode(),
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

// String satisfies the fuzzy.Source interface.
func (idx *Index) String(i int) string {
	s := lo.PickByKeys(
		idx.Data[i],
		idx.SearchableFields(),
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
