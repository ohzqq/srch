package index

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/ohzqq/srch/param"
)

const hareTestPath = `/home/mxb/code/srch/testdata/hare`
const hareTestURL = `file://home/mxb/code/srch/testdata/hare`
const hareTestQuery = `?url=file://home/mxb/code/srch/testdata/hare/`

var defTbls = []string{"_settings"}

type params struct {
	query string
	Cfg   *param.Cfg
}

type test struct {
	query string
	Cfg   *param.Cfg
}

func (p test) str() string {
	return p.query
}

func (p test) vals() url.Values {
	v, _ := url.ParseQuery(strings.TrimPrefix(p.query, "?"))
	return v
}

func (p test) url() *url.URL {
	u, _ := url.Parse(p.query)
	return u
}

func (p test) slice(got, want []string) error {
	if !slices.Equal(got, want) {
		return p.err(got, want)
	}
	return nil
}

func (p test) err(got, want any) error {
	return fmt.Errorf("query %v\ngot %#v, wanted %#v\n", p.str(), got, want)
}

func (t test) msg(msg any) error {
	return fmt.Errorf("query: %v\nerror: %#v\n", t.str(), msg)
}

func sliceTest(num, field any, got, want []string) error {
	if !slices.Equal(got, want) {
		return paramTestMsg(num, field, got, want)
	}
	return nil
}

func paramTestMsg(num, field, got, want any) error {
	return fmt.Errorf("test %v, field %s\ngot %#v, wanted %#v\n", num, field, got, want)
}
