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
	search      SearchFunc
	Fields      []*Field `json:"fields"`
	Query       Query    `json:"filters"`
	Identifier  string   `json:"identifier"`
	interactive bool
	indexed     bool
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

func (idx *Index) IndexData(data []map[string]any) *Results {
	idx.buildIndex(data)
	res := NewResults(data)
	if idx.HasFacets() {
		res.SetFacets(idx.Facets())
	}
	return res
}

func (idx *Index) FullTextSearch(src DataSrc, q string) *Results {
	res := idx.IndexData(src())

	if q == "" {
		return res
	}
	r := fullTextSearch(res.Data, idx.TextFields(), q)
	return idx.IndexData(r)
}

func (idx *Index) indexFacets(data []map[string]any) []*Field {
	return IndexData(data, idx.Facets(), idx.Identifier)
}

func (idx *Index) indexText(data []map[string]any) []*Field {
	return IndexData(data, idx.TextFields(), idx.Identifier)
}

func (idx *Index) buildIndex(data []map[string]any) *Index {
	idx.Fields = IndexData(data, idx.Fields, idx.Identifier)
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

	idx := make([]*Field, len(fields))
	for i, f := range fields {
		idx[i] = CopyField(f)
	}

	for _, d := range data {
		id := cast.ToUint32(d[id])
		for i, f := range fields {
			if val, ok := d[f.Attribute]; ok {
				idx[i].Add(val, id)
			}
		}
	}

	return idx
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
//func (idx *Index) FilterFacets(q any) *Index {
//  filters, err := NewQuery(q)
//  if err != nil {
//    log.Fatal(err)
//  }

//  idx.Query = filters
//  return FilterIndex(idx)
//}

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
