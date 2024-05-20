package srch

import (
	"net/url"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
)

type ClientCfg struct {
	tbl *hare.Table

	DB    string `json:"-" mapstructure:"path" qs:"db"`
	Index string `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	UID   string `query:"uid,omitempty" json:"uid,omitempty" mapstructure:"uid" qs:"uid"`
}

func NewClientCfg() *ClientCfg {
	return &ClientCfg{
		Index: "default",
	}
}

func (client *ClientCfg) Decode(v url.Values) error {
	err := sp.Decode(v, client)
	if err != nil {
		return err
	}
	return nil
}

func (client *ClientCfg) Encode() (url.Values, error) {
	v, err := sp.Encode(client)
	if err != nil {
		return nil, err
	}
	return v, nil
}
