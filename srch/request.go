package srch

import "github.com/ohzqq/srch/param"

type Request struct {
	*param.Params
}

func NewRequest(params string) (*Request, error) {
	p, err := param.Parse(params)
	if err != nil {
		return nil, err
	}

	req := &Request{
		Params: p,
	}

	return req, nil
}
