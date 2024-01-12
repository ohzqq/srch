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

func CfgFieldsFromValues(idx *Index, cfg url.Values) *Index {
	if cfg.Has("field") {
		for _, f := range cfg["field"] {
			idx.AddField(NewTextField(f))
		}
	}
	if cfg.Has("or") {
		for _, f := range cfg["or"] {
			idx.AddField(NewField(f, OrFacet))
		}
	}
	if cfg.Has("and") {
		for _, f := range cfg["and"] {
			idx.AddField(NewField(f, AndFacet))
		}
	}
	return idx
}
