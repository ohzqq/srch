package srch

import (
	"bytes"
	"log"
	"os"
)

type Opt func(*Index)

func Interactive(s *Index) {
	s.interactive = true
}

func WithCfgFile(file string) Opt {
	return func(idx *Index) {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = idx.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WithFields(fields []string) Opt {
	return func(idx *Index) {
		idx.SearchableFields = fields
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

// DataString initializes an index from a json formatted string.
func DataString(d string) Opt {
	return func(idx *Index) {
		buf := bytes.NewBufferString(d)
		err := idx.Decode(buf)
		if err != nil {
			log.Fatal(err)
		}

		if len(idx.Data) > 0 {
			idx.CollectItems()
		}
	}
}

func DataSlice(data []any) Opt {
	return func(idx *Index) {
		idx.Data = data
		idx.CollectItems()
	}
}

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
		idx.Data = data
		idx.CollectItems()
	}
}
