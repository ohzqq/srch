package param

import (
	"net/url"
)

type Cfg struct {
	FullText string   `query:"fullText,omitempty" json:"fullText,omitempty"`
	DataDir  string   `query:"dataDir,omitempty" json:"dataDir,omitempty"`
	DataFile []string `query:"dataFile,omitempty" json:"dataFile,omitempty"`
	UID      string   `query:"uid,omitempty" json:"uid,omitempty"`
}

func NewCfg() *Cfg {
	return &Cfg{}
}

func (s *Cfg) Parse(q string) error {
	vals, err := url.ParseQuery(q)
	if err != nil {
		return err
	}
	return s.Set(vals)
}

func (s *Cfg) Set(v url.Values) error {
	for _, key := range paramsCfg {
		switch key {
		case DataDir:
			s.DataDir = v.Get(key)
		case DataFile:
			s.DataFile = GetQueryStringSlice(key, v)
		case FullText:
			s.FullText = v.Get(key)
		case UID:
			s.UID = v.Get(key)
		}
		v.Del(key)
	}
	return nil
}

func (p Cfg) IsFullText() bool {
	return p.FullText != ""
}

func (p Cfg) HasData() bool {
	return len(p.DataFile) > 0 ||
		p.DataDir != ""
}

func (p Cfg) GetDataFiles() []string {
	var data []string
	switch {
	case p.DataDir != "":
		data = append(data, p.DataDir)
	case len(p.DataFile) > 0:
		data = append(data, p.DataFile...)
	}
	return data
}
