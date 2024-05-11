package idx

import (
	"net/url"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/db"
)

type Src interface {
	Find(...int) []any
}

type Tbl struct {
	*db.Table
}

type Mem struct {
	*hare.Database
}

type Net struct {
	*url.URL
}
