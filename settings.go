package srch

import (
	"net/url"
)

const (
	SrchAttr  = `searchableAttributes`
	FacetAttr = `attributesForFaceting`
)

type Settings struct {
	SearchableAttributes  []string `json:"searchableAttributes"`
	AttributesForFaceting []string `json:"attributesForFaceting"`
	TextAnalyzer          string   `json:"textAnalyzer"`
	SortFacetValuesBy     string   `json:"sortFacetValuesBy"`
}

func NewSettings(query any) *Settings {
	q := NewQuery(ParseQuery(query))

	s, err := q.Settings()
	if err != nil {
		return defaultSettings()
	}
	return s
}

func defaultSettings() *Settings {
	return &Settings{
		SearchableAttributes: []string{"title"},
		TextAnalyzer:         Fuzzy,
		SortFacetValuesBy:    "count",
	}
}

func GetAnalyzer(q url.Values) string {
	if q.Has("full_text") {
		return Text
	}
	return Fuzzy
}

func (s *Settings) Fields() []*Field {
	var fields []*Field
	fields = append(fields, s.TextFields()...)
	fields = append(fields, s.Facets()...)
	return fields
}

func (s *Settings) setValues(v url.Values) *Settings {
	q := NewQuery(v)
	s.SearchableAttributes = q.SrchAttr()

	s.AttributesForFaceting = q.FacetAttr()
	s.TextAnalyzer = q.Analyzer()
	return s
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
