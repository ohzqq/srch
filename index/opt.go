package index

type Opt func(*Idx) error

func WithRam(idx *Idx) error {
	return idx.memDB()
}
