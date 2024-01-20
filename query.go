package srch

import (
	"net/url"
	"strings"

	"github.com/spf13/cast"
)

type Query struct {
	params url.Values
}

func NewQuery(queries ...any) url.Values {
	q := make(url.Values)
	for _, query := range queries {
		vals, err := ParseValues(query)
		if err != nil {
			continue
		}
		for k, val := range vals {
			for _, v := range val {
				q.Add(k, v)
			}
		}
	}
	return q
}

func newQuery(q url.Values) *Query {
	return &Query{
		params: q,
	}
}

func (q Query) Data() ([]map[string]any, error) {
	return GetDataFromQuery(&q.params)
}

func (q Query) SrchAttr() []string {
	return GetQueryStringSlice(SrchAttr, q.params)
}

func (q Query) FacetAttr() []string {
	return GetQueryStringSlice(FacetAttr, q.params)
}

func (q Query) Analyzer() string {
	return GetAnalyzer(q.params)
}

func (q Query) Settings() *Settings {
	s := defaultSettings()
	return s.setValues(q.params)
}

func GetDataFromQuery(q *url.Values) ([]map[string]any, error) {
	var data []map[string]any
	var err error
	switch {
	case q.Has("data_file"):
		qu := *q
		data, err = FileSrc(qu["data_file"]...)
		q.Del("data_file")
	case q.Has("data_dir"):
		data, err = DirSrc(q.Get("data_dir"))
		q.Del("data_dir")
	}
	return data, err
}

func GetQueryStringSlice(key string, q url.Values) []string {
	var vals []string
	if q.Has(key) {
		for _, val := range q[key] {
			if val == "" {
				break
			}
			for _, v := range strings.Split(val, ",") {
				vals = append(vals, v)
			}
		}
	}
	if key == SrchAttr {
		switch len(vals) {
		case 0:
			return []string{"title"}
		case 1:
			if vals[0] == "" {
				return []string{"title"}
			}
		}
	}
	return vals
}

func ParseFieldsFromValues(cfg url.Values) []*Field {
	var fields []*Field
	if cfg.Has("field") {
		for _, f := range cfg["field"] {
			ft := Fuzzy
			if cfg.Has("full_text") {
				ft = Text
			}
			fields = append(fields, NewField(f, ft))
		}
	}
	if cfg.Has("or") {
		for _, f := range cfg["or"] {
			fields = append(fields, NewField(f, OrFacet))
		}
	}
	if cfg.Has("and") {
		for _, f := range cfg["and"] {
			fields = append(fields, NewField(f, AndFacet))
		}
	}
	return fields
}

// ParseValues takes an interface{} and returns a url.Values.
func ParseValues(f any) (url.Values, error) {
	filters := make(url.Values)
	var err error
	switch val := f.(type) {
	case url.Values:
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

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, err
}

// ParseQueryBytes parses a byte slice to url.Values.
func ParseQueryBytes(val []byte) (url.Values, error) {
	filters, err := cast.ToStringMapStringSliceE(string(val))
	if err != nil {
		return nil, err
	}
	return url.Values(filters), err
}
