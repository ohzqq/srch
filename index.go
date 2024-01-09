package srch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Lshortfile)
	viper.SetDefault("workers", 1)
}

// Index is a structure for facets and data.
type Index struct {
	src         DataSrc
	data        []map[string]any
	search      SearchFunc
	Fields      []*Field `json:"fields"`
	Query       Query    `json:"filters"`
	Identifier  string   `json:"identifier"`
	interactive bool
	fuzzy       bool
}

// New initializes an *Index with defaults: SearchableFields are
// []string{"title"}.
func New(src DataSrc, opts ...Opt) *Index {
	idx := &Index{
		data:       src(),
		src:        src,
		Identifier: "id",
		Fields:     []*Field{NewTextField("title")},
	}

	for _, opt := range opts {
		opt(idx)
	}

	idx.BuildIndex()

	return idx
}

// CopyIndex copies an index's config.
func CopyIndex(idx *Index, data []map[string]any) *Index {
	n := New(SliceSrc(data), WithCfg(idx.GetConfig()))
	n.Query = idx.Query
	n.search = idx.search
	n.interactive = idx.interactive
	return n
}

func (idx *Index) Data() []map[string]any {
	return idx.src()
}

func (idx *Index) BuildIndex() *Index {
	for _, d := range idx.Data() {
		id := cast.ToUint32(d[idx.Identifier])
		for _, f := range idx.Fields {
			if val, ok := d[f.Attribute]; ok {
				f.Add(val, id)
			}
		}
	}
	return idx
}

func DefaultIndex() *Index {
	return &Index{
		Identifier: "id",
		Fields:     []*Field{NewTextField("title")},
	}
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

func BuildIndex(data []map[string]any, opts ...Opt) *Index {
	idx := DefaultIndex()
	idx.data = data
	for _, opt := range opts {
		opt(idx)
	}
	idx.Fields = IndexData(data, idx.Fields, idx.Identifier)
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

// Filter idx.Data and re-calculate facets.
func (idx *Index) Filter(q any) *Index {
	filters, err := NewQuery(q)
	if err != nil {
		log.Fatal(err)
	}

	idx.Query = filters
	return Filter(idx)
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

// CfgIndex configures an *Index.
func CfgIndex(idx *Index, cfg any) {
	switch val := cfg.(type) {
	case []byte:
		err := CfgIndexFromBytes(idx, val)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
		return
	case string:
		if exist(val) {
			err := CfgIndexFromFile(idx, val)
			if err != nil {
				log.Printf("cfg error: %v, using defaults\n", err)
			}
			return
		} else {
			err := CfgIndexFromBytes(idx, []byte(val))
			if err != nil {
				log.Printf("cfg error: %v, using defaults\n", err)
			}
			return
		}
	case map[string]any:
		err := CfgIndexFromMap(idx, val)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
	}
}

// CfgIndexFromFile initializes an index from files.
func CfgIndexFromFile(idx *Index, cfg string) error {
	f, err := os.Open(cfg)
	if err != nil {
		return err
	}
	defer f.Close()

	err = idx.Decode(f)
	if err != nil {
		return err
	}

	return nil
}

// CfgIndexFromBytes initializes an index from a json formatted string.
func CfgIndexFromBytes(idx *Index, d []byte) error {
	err := idx.Decode(bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

// CfgIndexFromMap initalizes an index from a map[string]any.
func CfgIndexFromMap(idx *Index, d map[string]any) error {
	err := mapstructure.Decode(d, idx)
	if err != nil {
		return err
	}
	return nil
}

//func parseFacetMap(f any) map[string]*Facet {
//  facets := make(map[string]*Facet)
//  for name, agg := range cast.ToStringMap(f) {
//    facet := NewFacet(name)
//    err := mapstructure.Decode(agg, facet)
//    if err != nil {
//      log.Fatal(err)
//    }
//    facets[name] = facet
//  }
//  return facets
//}

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
