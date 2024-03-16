package param

import (
	"net/url"
)

type SrchCfg struct {
	BlvPath  string   `query:"fullText,omitempty" json:"fullText,omitempty"`
	DataDir  string   `query:"dataDir,omitempty" json:"dataDir,omitempty"`
	DataFile []string `query:"dataFile,omitempty" json:"dataFile,omitempty"`
	UID      string   `query:"uid,omitempty" json:"uid,omitempty"`
}

func NewCfg() *SrchCfg {
	return &SrchCfg{}
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
	return nil
}

func (s *SrchCfg) Has(key string) bool {
	switch key {
	case DataDir:
		return s.DataDir != ""
	case DataFile:
		return len(s.DataFile) > 0
	case FullText:
		return s.BlvPath != ""
	case UID:
		return s.UID != ""
	}
	return false
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
