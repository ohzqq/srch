package srch

import (
	"errors"
	"net/url"
	"strings"
)

type QueryParam interface {
	Decode(url.Values) error
	Encode() (url.Values, error)
}

func Decode(q any, p QueryParam) error {
	v, err := decodeQ(q)
	if err != nil {
		return err
	}
	return p.Decode(v)
}

func decodeQ(q any) (url.Values, error) {
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
