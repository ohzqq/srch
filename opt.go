package srch

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
)

type Opt func(*Index) Opt

func CfgFile(file string) Opt {
	return func(idx *Index) Opt {
		err := CfgIndexFromFile(idx, file)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
		return CfgFile(file)
	}
}

func CfgMap(m map[string]any) Opt {
	return func(idx *Index) Opt {
		err := CfgIndexFromMap(idx, m)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
		return CfgMap(m)
	}
}

func CfgString(cfg string) Opt {
	return CfgBytes([]byte(cfg))
}

func CfgBytes(cfg []byte) Opt {
	return func(idx *Index) Opt {
		err := CfgIndexFromBytes(idx, cfg)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
		return CfgBytes(cfg)
	}
}

func ReadCfg(r io.Reader) Opt {
	return func(idx *Index) Opt {
		err := idx.Decode(r)
		if err != nil {
			log.Printf("cfg error: %v, using defaults\n", err)
		}
		return ReadCfg(r)
	}
}

func WithCfg(c any) Opt {
	return func(idx *Index) Opt {
		CfgIndex(idx, c)
		return WithCfg(c)
	}
}

func WithFields(fields []*Field) Opt {
	return func(idx *Index) Opt {
		idx.AddField(fields...)
		return WithFields(fields)
	}
}

func WithFacets(fields []string) Opt {
	return func(idx *Index) Opt {
		idx.AddField(NewFacets(fields)...)
		return WithFacets(fields)
	}
}

func WithTextFields(fields []string) Opt {
	return func(idx *Index) Opt {
		idx.AddField(NewTextFields(fields)...)
		return WithTextFields(fields)
	}
}

func WithSearch(s SearchFunc) Opt {
	return func(idx *Index) Opt {
		search := idx.search
		idx.search = s
		return WithSearch(search)
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
