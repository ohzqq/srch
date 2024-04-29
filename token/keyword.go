package token

import "strings"

type Keyword struct {
	og []string
}

func (k *Keyword) Tokenize(og ...string) []string {
	k.og = og
	toks := make([]string, len(og))
	for i, t := range og {
		toks[i] = strings.ToLower(t)
	}
	return toks
}
