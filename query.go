package srch

import (
	"net/url"
)

type Query url.Values

// NewQuery takes an interface{} and returns a Query.
func NewQuery(f any) (Query, error) {
	val, err := ParseValues(f)
	return Query(val), err
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (Query, error) {
	f, err := FilterString(val)
	return Query(f), err
}

// ParseQueryBytes parses a byte slice to url.Values.
func ParseQueryBytes(val []byte) (Query, error) {
	f, err := FilterBytes(val)
	return Query(f), err
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
