package srch

import (
	"mime"
	"net/url"
	"path/filepath"
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
}

const (
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
)

type Data struct {
	*url.URL
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
