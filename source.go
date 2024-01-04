package srch

import "github.com/sahilm/fuzzy"

type Source interface {
	Items() []any
	fuzzy.Source
}
