package db

import (
	"github.com/ohzqq/hare"
)

type Net struct {
	*hare.Database
	name string
}
