package srch

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"os"
	"slices"
	"strings"

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
	search SearchFunc
	Fields []*Field         `json:"fields"`
	Data   []map[string]any `json:"data"`
	Query  url.Values       `json:"query"`
}

type SearchFunc func(string) []map[string]any

func New(q any, srch ...SearchFunc) *Index {
	idx := &Index{}
	idx.ParseQuery(q)

	if len(srch) > 0 {
		idx.search = srch[0]
	}

	switch {
	case idx.Query.Has("q"):
		return idx.Search(idx.Query.Get("q"))
	}

	return idx
}

func (idx *Index) Index(src []map[string]any) *Index {
	if len(idx.Fields) < 1 {
		idx.Fields = []*Field{NewField("title", Text)}
	}
	idx.Data = src
	if idx.Query.Has("sort_by") {
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

func (idx *Index) ParseQuery(q any) *Index {
	if q == nil {
		return idx.SetQuery(make(url.Values))
	}
	return idx.SetQuery(NewQuery(q))
}

func (idx *Index) SetQuery(q url.Values) *Index {
	idx.Query = q
	idx.AddField(ParseFieldsFromValues(idx.Query)...)

	data, err := GetDataFromQuery(&idx.Query)
	if err == nil {
		return idx.Index(data)
	}

	return idx
}

func (idx *Index) Sort() {
	sortDataByField(idx.Data, idx.Query.Get("sort_by"))
	if idx.Query.Has("order") {
		if idx.Query.Get("order") == "desc" {
			slices.Reverse(idx.Data)
		}
	}
}

func (idx *Index) Search(q string) *Index {
	if idx.search == nil {
		idx.search = FullTextSrchFunc(idx.Data, idx.TextFields())
	}
	idx.Query.Set("q", q)
	data := idx.search(q)
	res := idx.Copy().Index(data)

	if idx.HasFilters() {
		return idx.Filter(idx.Filters())
	}

	return res
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

func (idx *Index) Copy() *Index {
	return New(idx.Query)
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
	res := []string{
		"and",
		"or",
		"field",
		"q",
		"sort_by",
		"order",
		"data_file",
		"data_dir",
	}
	return lo.OmitByKeys(idx.Query, res)
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
		idx.ParseQuery(q)
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
		"query":  idx.Query.Encode(),
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
