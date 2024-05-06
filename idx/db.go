package idx

import (
	"github.com/ohzqq/srch/db"
	"github.com/ohzqq/srch/param"
)

func WithURL(uri string) db.Opt {
	return func(db *db.DB) error {
		return nil
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
	db, err := db.New(db.WithRam)
	if err != nil {
		return nil, err
	}

	if !db.TableExists(params.IndexName) {
		err = db.CreateTable(params.IndexName)
		if err != nil {
			return nil, err
		}

		m := NewMappingFromParams(params)
		_, err = db.CfgTable(params.IndexName, m)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
