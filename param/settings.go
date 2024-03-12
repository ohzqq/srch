package param

import (
	"mime"
	"net/url"
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
}

type Settings struct {
	FullText     string   `query:"fullText,omitempty" json:"fullText,omitempty"`
	SrchAttr     []string `query:"searchableAttributes,omitempty" json:"searchableAttributes,omitempty"`
	FacetAttr    []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty"`
	SortAttr     []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty"`
	DataDir      string   `query:"dataDir,omitempty" json:"dataDir,omitempty"`
	DataFile     []string `query:"dataFile,omitempty" json:"dataFile,omitempty"`
	DefaultField string   `query:"defaultField,omitempty" json:"defaultField,omitempty"`
	UID          string   `query:"uid,omitempty" json:"uid,omitempty"`

	params url.Values
}

func NewSettings() *Settings {
	return &Settings{
		params:   make(url.Values),
		SrchAttr: []string{DefaultField},
	}
}

func (s *Settings) Parse(v url.Values) error {
	for _, key := range paramsSettings {
		switch key {
		case SrchAttr:
			s.SrchAttr = parseSrchAttr(v)
		case FacetAttr:
			s.FacetAttr = parseFacetAttr(v)
		case SortAttr:
			s.SortAttr = GetQueryStringSlice(key, v)
		case DataDir:
			s.DataDir = v.Get(key)
		case DataFile:
			s.DataFile = GetQueryStringSlice(key, v)
		case DefaultField:
			s.DefaultField = v.Get(key)
		case FullText:
			s.FullText = v.Get(key)
		case UID:
			s.UID = v.Get(key)
		}
		v.Del(key)
	}
	return nil
}

func (p Settings) HasData() bool {
	return len(p.DataFile) > 0 ||
		p.DataDir != ""
}

func (p Settings) GetDataFiles() []string {
	var data []string
	switch {
	case p.DataDir != "":
		data = append(data, p.DataDir)
	case len(p.DataFile) > 0:
		data = append(data, p.DataFile...)
	}
	return data
}
