package srch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Lshortfile)
	viper.SetDefault("workers", 1)
}

// Index is a structure for facets and data.
type Index struct {
	Data        []map[string]any
	search      SearchFunc
	Fields      []*Field `json:"fields"`
	Query       Query    `json:"filters"`
	Identifier  string   `json:"identifier"`
	interactive bool
}

func New(opts ...Opt) *Index {
	idx := &Index{
		Identifier: "id",
	}
	for _, opt := range opts {
		opt(idx)
	}
	if len(idx.Fields) < 1 {
		idx.Fields = []*Field{NewTextField("title")}
	}
	return idx
}

func (idx *Index) Filter(src DataSrc, q string) (*Results, error) {
	data := src()
	idx.indexFacets(data)

	vals, err := ParseValues(q)
	if err != nil {
		return nil, err
	}
	items := Filter(data, idx.Facets(), vals)
	res := NewResults(items, idx.Facets()...)

	return res, nil
}

func (idx *Index) indexFacets(data []map[string]any) {
	idx.Fields = IndexData(data, idx.Facets(), idx.Identifier)
}

// NewIndex initializes an *Index with defaults: SearchableFields are
// []string{"title"}.
func NewIndex(src DataSrc, opts ...Opt) *Index {
	idx := New(opts...)

	//idx.BuildIndex()

	return idx
}

// CopyIndex copies an index's config.
func CopyIndex(idx *Index, data []map[string]any) *Index {
	n := NewIndex(SliceSrc(data), WithCfg(idx.GetConfig()))
	n.Query = idx.Query
	n.search = idx.search
	n.interactive = idx.interactive
	return n
}

func (idx *Index) GetData() []map[string]any {
	return idx.Data
}

func (idx *Index) BuildIndex() *Index {
	idx.Fields = IndexData(idx.Data, idx.Fields, idx.Identifier)
	return idx
}

func (idx *Index) AddField(fields ...*Field) *Index {
	idx.Fields = append(idx.Fields, fields...)
	return idx
}

func IndexData(data []map[string]any, fields []*Field, ident ...string) []*Field {
	id := "id"
	if len(ident) > 0 {
		id = ident[0]
	}
	for _, d := range data {
		id := cast.ToUint32(d[id])
		for _, f := range fields {
			if val, ok := d[f.Attribute]; ok {
				f.Add(val, id)
			}
		}
	}
	return fields
}

func IndexFacets(data []map[string]any, facets []string, ident ...string) []*Field {
	fields := NewFacets(facets)
	return IndexData(data, fields, ident...)
}

func IndexText(data []map[string]any, text []string, ident ...string) []*Field {
	fields := NewTextFields(text)
	return IndexData(data, fields, ident...)
}

func (idx *Index) GetField(attr string) (*Field, error) {
	for _, f := range idx.Fields {
		if f.Attribute == attr {
			return f, nil
		}
	}
	return nil, errors.New("no such field")
}

// Filter idx.Data and re-calculate facets.
func (idx *Index) FilterFacets(q any) *Index {
	filters, err := NewQuery(q)
	if err != nil {
		log.Fatal(err)
	}

	idx.Query = filters
	return FilterIndex(idx)
}

func (idx *Index) Facets() []*Field {
	return FilterFacets(idx.Fields)
}

func (idx *Index) TextFields() []*Field {
	return FilterTextFields(idx.Fields)
}

func (idx *Index) SearchableFields() []string {
	return SearchableFields(idx.Fields)
}

// GetConfig returns a map of the Index's config.
func (idx *Index) GetConfig() map[string]any {
	var facets []map[string]any
	for _, f := range idx.Fields {
		facets = append(facets, f.GetConfig())
	}
	return map[string]any{
		"fields":           facets,
		"searchableFields": idx.SearchableFields(),
	}
}

// HasFacets returns true if facets are configured.
func (idx *Index) HasFacets() bool {
	return len(idx.Facets()) > 0
}

// Decode unmarshals json from an io.Reader.
func (idx *Index) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(idx)
	if err != nil {
		return err
	}
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
