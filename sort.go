package srch

import (
	"slices"
	"strings"

	"github.com/spf13/cast"
)

type Sort struct {
	Field string
	Type  string
}

func NewSort(val string) *Sort {
	s := &Sort{}
	s.Field, s.Type = ParseSort(val)
	return s
}

func (s *Sort) Sort(data []map[string]any) []map[string]any {

	switch s.Type {
	case "int":
		s.SortByInt(data)
	case "string":
		s.SortByStr(data)
	}

	return data
}

func (s *Sort) SortByStr(data []map[string]any) {
	fn := func(a map[string]any, b map[string]any) int {
		var x, y string
		if v, ok := a[s.Field]; ok {
			x = cast.ToString(v)
		}
		if v, ok := b[s.Field]; ok {
			y = cast.ToString(v)
		}
		switch {
		case x > y:
			return 1
		case x == y:
			return 0
		default:
			return -1
		}
	}

	slices.SortStableFunc(data, fn)
}

func (s *Sort) SortByInt(data []map[string]any) {
	fn := func(a map[string]any, b map[string]any) int {
		var x, y int
		if v, ok := a[s.Field]; ok {
			x = cast.ToInt(v)
		}
		if v, ok := b[s.Field]; ok {
			y = cast.ToInt(v)
		}
		switch {
		case x > y:
			return 1
		case x == y:
			return 0
		default:
			return -1
		}
	}

	slices.SortStableFunc(data, fn)
}

func ParseSort(attr string) (string, string) {
	var by string
	t := "string"
	i := 0
	for attr != "" {
		var a string
		a, attr, _ = strings.Cut(attr, ":")
		if a == "" {
			continue
		}
		switch i {
		case 0:
			by = a
		case 1:
			t = a
		}
		i++
	}
	return by, t
}
