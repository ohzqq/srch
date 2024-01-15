package srch

import (
	"net/url"

	"github.com/spf13/cast"
)

func NewQuery(queries ...string) url.Values {
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

func GetDataFromQuery(q *url.Values) ([]map[string]any, error) {
	var data []map[string]any
	var err error
	switch {
	case q.Has("data_file"):
		data, err = dataFromFile(q.Get("data_file"))
		q.Del("data_file")
	case q.Has("data_dir"):
		data, err = DirSrc(q.Get("data_dir"))
		q.Del("data_dir")
	}
	return data, err
}

func ParseFieldsFromValues(cfg url.Values) []*Field {
	var fields []*Field
	if cfg.Has("field") {
		for _, f := range cfg["field"] {
			fields = append(fields, NewTextField(f))
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

func CfgFieldsFromValues(idx *Index, cfg url.Values) *Index {
	idx.AddField(ParseFieldsFromValues(cfg)...)
	return idx
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
