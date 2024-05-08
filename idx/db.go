package idx

import (
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/param"
)

func WithURL(uri string) db.Opt {
	fn := func(db *db.DB) error {
		return nil
	}
	return db.Opt{
		Name: "WithURL",
		Func: fn,
	}
}

type InitDB func() (*db.DB, error)

func NewDisk(params *param.Params) (*db.DB, error) {
	return db.New(db.NewDisk(params.Path))
}

func OpenDisk(params *param.Params) (*db.DB, error) {
	return db.New(db.WithDisk(params.Path))
}

func NewRam(params *param.Params) (*db.DB, error) {
	db, err := db.New(db.WithRam())
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
