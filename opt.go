package srch

import (
	"log"
	"os"
)

type Opt func(*Index)

func Interactive(s *Index) {
	s.interactive = true
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
