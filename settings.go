package srch

import (
	"net/url"
	"strings"
)

type Settings struct {
	SearchableAttributes  []string
	AttributesForFaceting []string
	TextAnalyzer          string
}

func NewSettings(query any) *Settings {
	settings := &Settings{
		SearchableAttributes: []string{"title"},
		TextAnalyzer:         Fuzzy,
	}

	q := NewQuery(query)

	settings.SearchableAttributes = GetQueryStringSlice("searchableAttributes", q)

	settings.AttributesForFaceting = GetQueryStringSlice("attributesForFaceting", q)
	settings.TextAnalyzer = GetAnalyzer(q)

	return settings
}

func GetAnalyzer(q url.Values) string {
	if q.Has("full_text") {
		return Text
	}
	return Fuzzy
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
	if key == "searchableAttributes" {
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

func (s *Settings) Fields() []*Field {
	var fields []*Field
	fields = append(fields, s.TextFields()...)
	fields = append(fields, s.Facets()...)
	return fields
}

func (s *Settings) TextFields() []*Field {
	fields := make([]*Field, len(s.SearchableAttributes))
	for i, attr := range s.SearchableAttributes {
		fields[i] = NewField(attr, s.TextAnalyzer)
	}
	return fields
}

func (s *Settings) Facets() []*Field {
	fields := make([]*Field, len(s.AttributesForFaceting))
	for i, attr := range s.AttributesForFaceting {
		fields[i] = NewField(attr, OrFacet)
	}
	return fields
}
