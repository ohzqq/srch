package srch

import (
	"net/url"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Query url.Values

// NewQuery takes an interface{} and returns a Query.
func NewQuery(f any) (Query, error) {
	val, err := ParseValues(f)
	return Query(val), err
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (Query, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return Query(q), err
}

// ParseQueryBytes parses a byte slice to url.Values.
func ParseQueryBytes(val []byte) (Query, error) {
	filters, err := cast.ToStringMapStringSliceE(string(val))
	if err != nil {
		return nil, err
	}
	return Query(filters), err
}

func (q Query) String() string {
	return q.Encode()
}

func (q Query) Values() url.Values {
	return url.Values(q)
}

func (q Query) Encode() string {
	return q.Values().Encode()
}

func (q Query) Set(k, v string) {
	q.Values().Set(k, v)
}

func (q Query) Keywords() []string {
	if q.Values().Has("q") {
		return q["q"]
	}
	return []string{}
}

func (q Query) Filters() url.Values {
	return lo.OmitByKeys(q, []string{"q"})
}

// ParseValues takes an interface{} and returns a url.Values.
func ParseValues(f any) (map[string][]string, error) {
	filters := make(map[string][]string)
	var err error
	switch val := f.(type) {
	case url.Values:
		return val, nil
	case Query:
		return val, nil
	case []byte:
		return ParseQueryBytes(val)
	case string:
		return ParseQueryString(val)
	default:
		filters, err = cast.ToStringMapStringSliceE(val)
		if err != nil {
			return nil, err
		}
	}
	return filters, nil
}
