package kw

import "github.com/spf13/cast"

type Field struct {
}

type Token struct {
}

func KeywordAnalyzer(val any) []*Token {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*Token, len(tokens))
	for i, token := range tokens {
		items[i] = NewToken(token)
		items[i].Value = normalizeText(token)
	}
	return items
}
