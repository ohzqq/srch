package token

import "strings"

func TokenizeKeywords(og []string) []string {
	toks := make([]string, len(og))
	for i, t := range og {
		toks[i] = strings.ToLower(t)
	}
	return toks
}
