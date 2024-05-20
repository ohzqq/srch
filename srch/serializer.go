package srch

import (
	"net/url"
	"strings"
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

func ParseQueryStrings(q []string) []string {
	var vals []string
	for _, val := range q {
		if val == "" {
			break
		}
		for _, v := range strings.Split(val, ",") {
			vals = append(vals, v)
		}
	}
	return vals
}
