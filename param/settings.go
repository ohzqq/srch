package param

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type Settings struct {
	FullText     string   `query:"fullText,omitempty" json:"fullText,omitempty" mapstructure:"fullText,omitempty"`
	SrchAttr     []string `query:"searchableAttributes,omitempty" json:"searchableAttributes,omitempty" mapstructure:"searchableAttributes,omitempty"`
	FacetAttr    []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributesForFaceting,omitempty"`
	SortAttr     []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortableAttributes,omitempty"`
	DataDir      string   `query:"dataDir,omitempty" json:"dataDir,omitempty" mapstructure:"dataDir,omitempty"`
	DataFile     []string `query:"dataFile,omitempty" json:"dataFile,omitempty" mapstructure:"dataFile,omitempty"`
	DefaultField string   `query:"defaultField,omitempty" json:"defaultField,omitempty" mapstructure:"defaultField,omitempty"`
	UID          string   `query:"uid,omitempty" json:"uid,omitempty" mapstructure:"uid,omitempty"`

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

func (p Settings) GetData() ([]map[string]any, error) {
	files, err := p.GetDataFiles()
	if err != nil {
		return nil, err
	}

	var data []map[string]any
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		err = DecodeData(f, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (p Settings) GetDataFiles() ([]string, error) {
	switch {
	case p.DataDir != "":
		d, err := filepath.Glob(filepath.Join(p.DataDir, "*.json"))
		if err != nil {
			return nil, err
		}
		return d, nil
	case len(p.DataFile) > 0:
		return p.DataFile, nil
	default:
		return []string{}, errors.New("no data")
	}
}

func (p *Settings) IsFullText() bool {
	return p.params.Has(FullText)
}

func (p *Settings) GetFullText() string {
	return p.params.Get(FullText)
}

func (p Settings) GetUID() string {
	return p.params.Get("uid")
}

// DecodeData decodes data from a io.Reader.
func DecodeData(r io.Reader, data *[]map[string]any) error {
	dec := json.NewDecoder(r)
	for {
		m := make(map[string]any)
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		*data = append(*data, m)
	}
	return nil
}
