package db

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type Record struct {
	DataID int      `json:"id,omitempty"`
	IdxID  int      `json:"_id"`
	Table  string   `json:"table"`
	Doc    *doc.Doc `json:"-"`
}

func (r *Record) SetID(id int) {
	r.IdxID = id
}

func (r *Record) GetID() int {
	return r.IdxID
}

func (r *Record) AfterFind(_ *hare.Database) error {
	return nil
}

//	func (db *DB) TableExists(name string) bool {
//	 if _, ok := db.Tables[name]; ok {
//	   return true
//	 }
//	 return false
//	}
