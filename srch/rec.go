package srch

import "github.com/ohzqq/hare"

type Item struct {
	ID   int            `json:"_id"`
	Data map[string]any `json:"data"`
}

func (r *Item) SetID(id int) {
	r.ID = id
}

func (r *Item) GetID() int {
	return r.ID
}

func (r *Item) AfterFind(_ *hare.Database) error {
	return nil
}
