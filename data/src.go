package data

import (
	"net/url"

	"github.com/ohzqq/srch/db"
)

type Src interface {
	Find(...int) []any
}

type Tbl struct {
	*db.Table
}

type Mem struct {
	*db.DB
}

type Net struct {
	*url.URL
}
