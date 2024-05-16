package param

import (
	"errors"
	"net/url"
	"strings"
)

type QueryParam interface {
	Decode(string) error
	Encode() (url.Values, error)
}

func Decode(q any, p QueryParam) error {
	switch v := q.(type) {
	case string:
		return p.Decode(v)
	case map[string][]string:
		return p.Decode(url.Values(v).Encode())
	case url.Values:
		return p.Decode(v.Encode())
	case *url.URL:
		return p.Decode(v.Query().Encode())
	default:
		return errors.New("param must be of type: string, map[string][]string, url.Values, *url.URL")
	}
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

func parseQuery(q string) url.Values {
	q = strings.TrimPrefix(q, "?")
	u, err := url.ParseQuery(q)
	if err != nil {
		return make(url.Values)
	}
	return u
}
