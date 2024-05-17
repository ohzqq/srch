package index

import (
	"net/url"
	"strings"
)

const hareTestPath = `/home/mxb/code/srch/testdata/hare`
const hareTestURL = `file://home/mxb/code/srch/testdata/hare`
const hareTestQuery = `?url=file://home/mxb/code/srch/testdata/hare/`

var defTbls = []string{"_settings"}

type params struct {
	query string
}

func (p params) str() string {
	return p.query
}

func (p params) vals() url.Values {
	v, _ := url.ParseQuery(strings.TrimPrefix(p.query, "?"))
	return v
}

func (p params) url() *url.URL {
	u, _ := url.Parse(p.query)
	return u
}
