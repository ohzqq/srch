package db

import "github.com/ohzqq/hare"

type Record struct {
	DataID int `json:"id,omitempty"`
	IdxID  int `json:"_id"`
	Data   any `json:"data"`
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
