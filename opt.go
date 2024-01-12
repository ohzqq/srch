package srch

import (
	"bytes"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Opt func(*Index)

func CfgFile(file string) Opt {
	return func(idx *Index) {
		err := CfgIndexFromFile(idx, file)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
	}
}

func CfgMap(m map[string]any) Opt {
	return func(idx *Index) {
		err := CfgIndexFromMap(idx, m)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
	}
}

func CfgString(cfg string) Opt {
	return CfgBytes([]byte(cfg))
}

func CfgBytes(cfg []byte) Opt {
	return func(idx *Index) {
		err := CfgIndexFromBytes(idx, cfg)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
	}
}

func ReadCfg(r io.Reader) Opt {
	return func(idx *Index) {
		err := idx.Decode(r)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
	}
}

func WithCfg(c any) Opt {
	return func(idx *Index) {
		CfgIndex(idx, c)
	}
}

func WithFields(fields []*Field) Opt {
	return func(idx *Index) {
		idx.AddField(fields...)
	}
}

func WithFacets(fields []string) Opt {
	return func(idx *Index) {
		idx.AddField(NewFacets(fields)...)
	}
}

func WithTextFields(fields []string) Opt {
	return func(idx *Index) {
		idx.AddField(NewTextFields(fields)...)
	}
}

func WithSearch(s SearchFunc) Opt {
	return func(idx *Index) {
		idx.search = s
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

func ParseCfgQuery(q string) (*Index, error) {
	v, err := url.ParseQuery(testValuesCfg)
	if err != nil {
		return New(), err
	}
	return CfgIndexFromValues(v)
}

func CfgIndexFromValues(cfg url.Values) (*Index, error) {
	idx := New()
	idx.Query = cfg
	CfgFieldsFromValues(idx, cfg)
	return idx, nil
}

func CfgFieldsFromValues(idx *Index, cfg url.Values) *Index {
	if cfg.Has("field") {
		for _, f := range cfg["field"] {
			idx.AddField(NewTextField(f))
		}
	}
	if cfg.Has("or") {
		for _, f := range cfg["or"] {
			idx.AddField(NewField(f, OrFacet))
		}
	}
	if cfg.Has("and") {
		for _, f := range cfg["and"] {
			idx.AddField(NewField(f, AndFacet))
		}
	}
	return idx
}

// ReadIdxCfg reads an *Index's config.
func ReadIdxCfg(r io.Reader) (*Index, error) {
	cfg := make(url.Values)
	err := yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return &Index{}, err
	}

	return CfgIndexFromValues(cfg)
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
