package srch

import (
	"errors"
	"net/url"
	"strings"
)

type Request struct {
	vals url.Values
}

func NewRequest(u any) (*Request, error) {
	v, err := ParseQuery(u)
	if err != nil {
		return nil, err
	}
	return &Request{vals: v}, nil
}

func (req *Request) Decode(u any) (*Client, error) {
	v, err := ParseQuery(u)
	if err != nil {
		return nil, err
	}

	cfg := NewCfg()
	err = cfg.Decode(v)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(v)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ParseQuery(q any) (url.Values, error) {
	switch v := q.(type) {
	case string:
		v = strings.TrimPrefix(v, "?")
		return url.ParseQuery(v)
	case map[string][]string:
		return url.Values(v), nil
	case url.Values:
		return v, nil
	case *url.URL:
		return v.Query(), nil
	default:
		return nil, errors.New("param must be of type: string, map[string][]string, url.Values, *url.URL")
	}
}
