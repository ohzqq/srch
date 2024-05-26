package srch

import (
	"mime"
	"net/url"
	"path/filepath"

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
	//switch u.Scheme {
	//case "file":
	//case "http", "https":
	//case "":
	//}
	return &Data{URL: u}
}

func (d *Data) ContentType() string {
	return mime.TypeByExtension(filepath.Ext(d.Path))
}
