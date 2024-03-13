package sep

import (
	"strings"
	"unicode"
)

type Func func(r rune) bool

func Split(str string, seps ...Func) []string {
	if len(seps) == 0 {
		return []string{str}
	}
	return strings.FieldsFunc(str, seps[0])
}

func Whitespace(r rune) bool { return unicode.IsSpace(r) }
func Comma(r rune) bool      { return r == ',' }
func Space(r rune) bool      { return r == ' ' }
func Tab(r rune) bool        { return r == '\t' }
func Newline(r rune) bool    { return r == '\r' || r == '\n' }
