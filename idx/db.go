package idx

import (
	"github.com/ohzqq/srch/client"
	"github.com/ohzqq/srch/param"
)

func WithURL(uri string) client.Opt {
	fn := func(client *client.Client) error {
		return nil
	}
	return client.Opt{
		Name: "WithURL",
		Func: fn,
	}
}

type InitDB func() (*client.Client, error)

func NewDisk(params *param.Params) (*client.Client, error) {
	return client.New(client.NewDisk(params.Path))
}

func OpenDisk(params *param.Params) (*client.Client, error) {
	return client.New(client.WithDisk(params.Path))
}

func NewRam(params *param.Params) (*client.Client, error) {
	db, err := client.New(client.WithRam())
	if err != nil {
		return nil, err
	}

	if !db.TableExists(params.IndexName) {
		err = db.CreateTable(params.IndexName)
		if err != nil {
			return nil, err
		}

		m := NewMappingFromParams(params)
		err = db.CfgTable(params.IndexName, m, params.UID)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
