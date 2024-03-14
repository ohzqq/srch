package param

import (
	"net/url"
)

type SrchCfg struct {
	BlvPath  string   `query:"fullText,omitempty" json:"fullText,omitempty"`
	DataDir  string   `query:"dataDir,omitempty" json:"dataDir,omitempty"`
	DataFile []string `query:"dataFile,omitempty" json:"dataFile,omitempty"`
	UID      string   `query:"uid,omitempty" json:"uid,omitempty"`
	*IndexSettings
}

func NewCfg() *SrchCfg {
	return &SrchCfg{
		IndexSettings: NewSettings(),
	}
}

func (s *SrchCfg) Parse(q string) error {
	vals, err := url.ParseQuery(q)
	if err != nil {
		return err
	}
	return s.Set(vals)
}

func (s *SrchCfg) Set(v url.Values) error {
	for _, key := range paramsCfg {
		switch key {
		case DataDir:
			s.DataDir = v.Get(key)
		case DataFile:
			s.DataFile = GetQueryStringSlice(key, v)
		case FullText:
			s.BlvPath = v.Get(key)
		case UID:
			s.UID = v.Get(key)
		}
		v.Del(key)
	}
	s.IndexSettings.Set(v)
	return nil
}

func (p SrchCfg) IsFullText() bool {
	return p.BlvPath != ""
}

func (p SrchCfg) HasData() bool {
	return len(p.DataFile) > 0 ||
		p.DataDir != ""
}

func (p SrchCfg) GetDataFiles() []string {
	var data []string
	switch {
	case p.DataDir != "":
		data = append(data, p.DataDir)
	case len(p.DataFile) > 0:
		data = append(data, p.DataFile...)
	}
	return data
}
