package srch

import "net/url"

type Query url.Values

func NewQuery(queries ...string) url.Values {
	q := make(url.Values)
	for _, query := range queries {
		vals, err := url.ParseQuery(query)
		if err != nil {
			break
		}
		for k, val := range vals {
			for _, v := range val {
				q.Add(k, v)
			}
		}
	}
	return q
}

func ParseCfgQuery(q string) (*Index, error) {
	v, err := url.ParseQuery(testValuesCfg)
	if err != nil {
		return New(), err
	}
	return CfgIndexFromValues(v)
}

func CfgIndexFromValues(cfg url.Values) (*Index, error) {
	idx := New()
	idx.Query = cfg
	CfgFieldsFromValues(idx, cfg)
	return idx, nil
}

func FieldsFromQuery(cfg url.Values) []*Field {
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
	idx.Fields = FieldsFromQuery(cfg)
	return idx
}
