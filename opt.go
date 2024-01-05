package srch

import (
	"bytes"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
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

func WithCfg(c any) Opt {
	return func(idx *Index) {
		switch val := c.(type) {
		case []byte:
			err := idx.Decode(bytes.NewBuffer(val))
			if err != nil {
				log.Fatal(err)
			}
			return
		case string:
			if exist(val) {
				f, err := os.Open(val)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				err = idx.Decode(f)
				if err != nil {
					log.Fatal(err)
				}

				return
			} else {
				err := idx.Decode(bytes.NewBufferString(val))
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		case map[string]any:
			err := mapstructure.Decode(val, idx)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func WithFacets(facets []*Facet) Opt {
	return func(idx *Index) {
		idx.Facets = facets
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
