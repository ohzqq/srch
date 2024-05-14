package index

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/param"
)

type Idx struct {
	*hare.Database
	*param.Params

	Cfg    *Cfg
	Tables map[string]int
}
