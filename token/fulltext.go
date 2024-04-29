package token

type Fulltext struct {
	og []string
}

func (f *Fulltext) Tokenize(og ...string) []string {
	f.og = og
	var toks []string
	for _, v := range og {
		tokens := Split(v)
		tokens = RemoveStopwords(tokens...)
		for _, t := range tokens {
			t = Stem(t)
			toks = append(toks, normalizeStr(t))
		}
	}
	return toks
}
