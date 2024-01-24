package txt

import "github.com/spf13/cast"

type Analyzer interface {
	Tokenize(any) []*Token
}

type Simple struct{}

func (s Simple) Tokenize(str any) []*Token {
	return []*Token{NewToken(cast.ToString(str))}
}
