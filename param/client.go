package param

import (
	"net/url"

	"github.com/ohzqq/sp"
)

type Client struct {
	*url.URL `json:"-"`

	DB    string `json:"-" mapstructure:"path" qs:"db"`
	Index string `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	UID   string `query:"uid,omitempty" json:"uid,omitempty" mapstructure:"uid" qs:"uid"`
}

func DefaultClient() *Client {
	return &Client{
		Index: "default",
		URL:   &url.URL{},
	}
}

func (client *Client) Decode(v url.Values) error {
	err := sp.Decode(v, client)
	if err != nil {
		return err
	}
	client.URL, err = parseURL(client.DB)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Encode() (url.Values, error) {
	v, err := sp.Encode(client)
	if err != nil {
		return nil, err
	}
	return v, nil
}
