package srch

import (
	"mime"
	"net/url"

	"github.com/ohzqq/hare"
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
	mime.AddExtensionType(".hare", "application/hare")
}

const (
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
	Hare   = `application/hare`
)

type Data struct {
	*url.URL

	db *hare.Database
}

func NewData(u *url.URL) *Data {
	data := &Data{
		URL: u,
	}
	return data
}
