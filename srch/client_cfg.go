package srch

import (
	"net/url"

	"github.com/ohzqq/hare"
)

type ClientCfg struct {
	*hare.Database
	*url.URL `json:"-"`

	tbl *hare.Table

	DB    string `json:"-" mapstructure:"path" qs:"db"`
	Index string `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	UID   string `query:"uid,omitempty" json:"uid,omitempty" mapstructure:"uid" qs:"uid"`
}

func NewClientCfg() *ClientCfg {
	return &ClientCfg{
		Index: "default",
		URL:   &url.URL{},
	}
}
