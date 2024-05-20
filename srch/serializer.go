package srch

import (
	"net/url"
)

type QueryParam interface {
	Decode(url.Values) error
	Encode() (url.Values, error)
}

func Decode(q any, p QueryParam) error {
	v, err := ParseQuery(q)
	if err != nil {
		return err
	}
	return p.Decode(v)
}

func Encode(p QueryParam) (url.Values, error) {
	return p.Encode()
}
