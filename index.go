package srch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

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
	search SearchFunc
	Fields []*Field         `json:"fields"`
	Data   []map[string]any `json:"data"`
	Query  url.Values       `json:"query"`
}

type SearchFunc func(string) []map[string]any

func OldNew(opts ...Opt) *Index {
	idx := &Index{
		Query: make(url.Values),
	}
	for _, opt := range opts {
		opt(idx)
	}
	return idx
}

func New(q ...string) *Index {
	idx := &Index{
		Query: NewQuery(q...),
	}

	idx.AddField(FieldsFromQuery(idx.Query)...)

	data, err := GetDataFromQuery(&idx.Query)
	if err == nil {
		idx.Index(data)
	}

	return idx
}

func (idx *Index) Index(src []map[string]any) *Index {
	if len(idx.Fields) < 1 {
		idx.Fields = []*Field{NewTextField("title")}
	}
	idx.Data = src
	idx.Fields = IndexData(idx.Data, idx.Fields)
	return idx
}

func (idx *Index) Facets() []*Facet {
	return FieldsToFacets(idx.FacetFields())
}

func (idx *Index) Search(q string) *Index {
	if idx.search == nil {
		idx.search = FullTextSrchFunc(idx.Data, idx.TextFields())
	}
	idx.Query.Set("q", q)
	res := idx.search(q)
	return idx.Copy().Index(res)
}

func (idx *Index) Filter(q string) *Index {
	vals, err := ParseValues(q)
	if err != nil {
		return idx
	}
	data := Filter(idx.Data, idx.FacetFields(), vals)
	return idx.Copy().Index(data)
}

func (idx *Index) AddField(fields ...*Field) *Index {
	idx.Fields = append(idx.Fields, fields...)
	return idx
}

func IndexData(data []map[string]any, fields []*Field) []*Field {
	for _, f := range fields {
		f.Items = make(map[string]*roaring.Bitmap)
	}

	for id, d := range data {
		for i, f := range fields {
			if val, ok := d[f.Attribute]; ok {
				fields[i].Add(val, uint32(id))
			}
		}
	}

	return fields
}

func (idx *Index) GetField(attr string) (*Field, error) {
	for _, f := range idx.Fields {
		if f.Attribute == attr {
			return f, nil
		}
	}
	return nil, errors.New("no such field")
}

func (idx *Index) GetFacet(attr string) (*Facet, error) {
	for _, f := range idx.Fields {
		if f.Attribute == attr {
			return NewFacet(f), nil
		}
	}
	return nil, errors.New("no such field")
}

func (idx *Index) String(i int) string {
	s := lo.PickByKeys(
		idx.Data[i],
		idx.SearchableFields(),
	)
	vals := cast.ToStringSlice(lo.Values(s))
	return strings.Join(vals, "\n")
}

func (idx *Index) Len() int {
	return len(idx.Data)
}

func (idx *Index) AddFieldsFromValues(cfg url.Values) *Index {
	return CfgFieldsFromValues(idx, cfg)
}

func (idx *Index) CfgString() string {
	return idx.Query.Encode()
}

func (idx *Index) FacetFields() []*Field {
	return FilterFacets(idx.Fields)
}

func (idx *Index) FilterByID(ids []int) *Index {
	data := FilterDataByID(idx.Data, ids)
	return idx.Copy().Index(data)
}

func (idx *Index) Copy() *Index {
	return OldNew(WithFields(idx.Fields))
}

func (idx *Index) TextFields() []*Field {
	return FilterTextFields(idx.Fields)
}

func (idx *Index) SearchableFields() []string {
	return SearchableFields(idx.Fields)
}

// HasFacets returns true if facets are configured.
func (idx *Index) HasFacets() bool {
	return len(idx.FacetFields()) > 0
}

// Decode unmarshals json from an io.Reader.
func (idx *Index) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(idx)
	//err := json.NewDecoder(r).Decode(&idx.Query)
	if err != nil {
		return err
	}
	//idx.AddFieldsFromValues(idx.Query)
	return nil
}

// Encode marshals json from an io.Writer.
func (idx *Index) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(idx)
}

// JSON marshals an Index to json.
func (idx *Index) JSON() []byte {
	var buf bytes.Buffer
	err := idx.Encode(&buf)
	if err != nil {
		return []byte("{}")
	}
	return buf.Bytes()
}

func (idx *Index) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"data":   idx.Data,
		"facets": idx.Facets(),
		"query":  idx.Query.Encode(),
	}
	return json.Marshal(res)
}

// Print writes Index json to stdout.
func (idx *Index) Print() {
	err := idx.Encode(os.Stdout)
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

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
