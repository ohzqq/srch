package srch

import "strings"

type Settings struct {
	SearchableAttributes  []string
	AttributesForFaceting []string
	TextAnalyzer          string
}

func NewSettings(q any) *Settings {
	settings := &Settings{
		SearchableAttributes: []string{"title"},
		TextAnalyzer:         Fuzzy,
	}

	v, err := ParseValues(q)
	if err != nil {
		return settings
	}

	if len(v) < 1 {
		return settings
	}

	for k, vals := range v {
		var attr []string

		switch len(vals) {
		case 0:
			break
		case 1:
			if vals[0] != "" {
				attr = strings.Split(vals[0], ",")
			}
		default:
			attr = vals
		}

		if len(attr) < 1 {
			break
		}

		switch k {
		case "searchableAttributes":
			settings.SearchableAttributes = attr
		case "attributesForFaceting":
			settings.AttributesForFaceting = attr
		case "textAnalyzer":
			settings.TextAnalyzer = attr[0]
		}
	}

	return settings
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
