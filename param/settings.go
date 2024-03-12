package param

import (
	"encoding/json"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

		ct := mime.TypeByExtension(filepath.Ext(file))

		err = DecodeData(f, ct, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (p Settings) GetDataFiles() ([]string, error) {
	var data []string
	switch {
	case p.DataDir != "":
		d, err := filepath.Glob(filepath.Join(p.DataDir, "*.json"))
		if err != nil {
			return nil, err
		}
		data = append(data, d...)
		nd, err := filepath.Glob(filepath.Join(p.DataDir, "*.ndjson"))
		if err != nil {
			return nil, err
		}
		data = append(data, nd...)
	case len(p.DataFile) > 0:
		return p.DataFile, nil
	}
	return data, nil
}

func DecodeData(r io.Reader, ct string, data *[]map[string]any) error {
	b, _, ok := strings.Cut(ct, ";")
	if ok {
		ct = b
	}
	switch ct {
	case "application/x-ndjson":
		return DecodeNDJSON(r, data)
	case "application/json":
		return DecodeJSON(r, data)
	}
	return nil
}

// DecodeNDJSON decodes data from a io.Reader.
func DecodeNDJSON(r io.Reader, data *[]map[string]any) error {
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

func DecodeJSON(r io.Reader, data *[]map[string]any) error {
	dec := json.NewDecoder(r)
	for {
		m := []map[string]any{}
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		*data = append(*data, m...)
	}
	return nil
}
