package srch

import (
	"bytes"
	"io"
	"log"
	"os"
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
		idx.Fields = fields
	}
}

func WithFacets(fields []string) Opt {
	return func(idx *Index) {
		for _, f := range fields {
			idx.Fields = append(idx.Fields, NewTaxonomyField(f))
		}
	}
}

func WithTextFields(fields []string) Opt {
	return func(idx *Index) {
		for _, f := range fields {
			idx.Fields = append(idx.Fields, NewTextField(f))
		}
	}
}

func WithSearch(s SearchFunc) Opt {
	return func(idx *Index) {
		idx.search = s
	}
}

func WithFuzzySearch() Opt {
	return func(idx *Index) {
		idx.fuzzy = true
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
		idx.src = SliceSrc(data)
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

		data, err := DecodeData(f)
		if err != nil {
			log.Fatal(err)
		}
		idx.src = SliceSrc(data)
	}
}
