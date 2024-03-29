package srch

import "strings"

type Sort struct {
	Field string
	Order string
	Type  string
}

func NewSort(val string) *Sort {
	s := &Sort{}
	s.Field, s.Order, s.Type = ParseSort(val)
	return s
}

func ParseSort(attr string) (string, string, string) {
	var by string
	order := "desc"
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
			order = a
		case 2:
			t = a
		}
		i++
	}
	return by, order, t
}
