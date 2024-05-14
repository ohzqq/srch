package index

import (
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/doc"
)

type Cfg struct {
	*hare.Table `json:"-"`

	ID       int         `json:"_id"`
	Name     string      `json:"name"`
	CustomID string      `json:"customID,omitempty"`
	Mapping  doc.Mapping `json:"mapping"`
}
