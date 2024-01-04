package srch

type Opt func(*Index)

func Interactive(s *Index) {
	s.interactive = true
}

func WithFields(fields []string) Opt {
	return func(idx *Index) {
		idx.SearchFields = fields
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
