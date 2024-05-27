package srch

import (
	"slices"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
)

type Item struct {
	ID   int            `json:"_id"`
	Data map[string]any `json:"data"`
}

func NewItem() *Item {
	return &Item{
		Data: make(map[string]any),
	}
}

func (i *Item) Idx(m Mapping) *Doc {
	doc := DefaultDoc()
	for ana, attrs := range m {
		for field, val := range i.Data {
			for _, attr := range attrs {
				if field == attr {
					if ana == analyzer.Simple && slices.Equal(attrs, []string{"*"}) {
						doc.AddField(ana, field, val)
					}
					doc.AddField(ana, field, val)
				}
			}
		}
	}
	return doc
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
