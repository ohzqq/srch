package param

import (
	"mime"
	"net/url"
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
}

type Settings struct {
	SrchAttr     []string `query:"searchableAttributes,omitempty" json:"searchableAttributes,omitempty"`
	FacetAttr    []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty"`
	SortAttr     []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty"`
	DefaultField string   `query:"defaultField,omitempty" json:"defaultField,omitempty"`
	UID          string   `query:"uid,omitempty" json:"uid,omitempty"`
}

func NewSettings() *Settings {
	return &Settings{
		SrchAttr: []string{DefaultField},
	}
}

func (s *Settings) Parse(q string) error {
	vals, err := url.ParseQuery(q)
	if err != nil {
		return err
	}
	return s.Set(vals)
}

func (s *Settings) Set(v url.Values) error {
	for _, key := range paramsSettings {
		switch key {
		case SrchAttr:
			s.SrchAttr = parseSrchAttr(v)
		case FacetAttr:
			s.FacetAttr = parseFacetAttr(v)
		case SortAttr:
			s.SortAttr = GetQueryStringSlice(key, v)
		case DefaultField:
			s.DefaultField = v.Get(key)
		case UID:
			s.UID = v.Get(key)
		}
		v.Del(key)
	}
	return nil
}

func parseSrchAttr(vals url.Values) []string {
	if !vals.Has(SrchAttr) {
		return []string{DefaultField}
	}
	vals[SrchAttr] = GetQueryStringSlice(SrchAttr, vals)
	if len(vals[SrchAttr]) < 1 {
		vals[SrchAttr] = []string{DefaultField}
	}
	return vals[SrchAttr]
}

func parseFacetAttr(vals url.Values) []string {
	if !vals.Has(Facets) {
		vals[Facets] = GetQueryStringSlice(FacetAttr, vals)
	}
	return vals[Facets]
}
