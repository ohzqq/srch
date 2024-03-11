package param

import "net/url"

type Params struct {
	Settings url.Values
	Search   url.Values
}

func NewParams() *Params {
	p := &Params{
		Settings: make(url.Values),
		Search:   make(url.Values),
	}
	return p
}
