package srch

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
)

type Opt func(*Index)

func Interactive(s *Index) {
	s.interactive = true
}

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

// DataString sets the *Index.Data from a json formatted string.
func DataString(d string) Opt {
	return func(idx *Index) {
		buf := bytes.NewBufferString(d)
		err := idx.Decode(buf)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// DataSlice sets the *Index.Data from a slice.
func DataSlice(data []map[string]any) Opt {
	return func(idx *Index) {
		//idx.Data = data
	}
}

// DataFile sets the *Index.Data from a json file.
func DataFile(cfg string) Opt {
	return func(idx *Index) {
		f, err := os.Open(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		//data, err := DecodeData(f)
		//if err != nil {
		//log.Fatal(err)
		//}
		//idx.Data = data
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
