package srch

import (
	"net/url"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
)

type ClientCfg struct {
	*hare.Table

	DB    string `json:"-" mapstructure:"path" qs:"db"`
	Index string `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	UID   string `query:"uid,omitempty" json:"uid,omitempty" mapstructure:"uid" qs:"uid"`
}

func NewClientCfg() *ClientCfg {
	return &ClientCfg{
		Index: "default",
	}
}

func (cfg *ClientCfg) Decode(v url.Values) error {
	err := sp.Decode(v, cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *ClientCfg) Encode() (url.Values, error) {
	v, err := sp.Encode(cfg)
	if err != nil {
		return nil, err
	}
	return v, nil
}
