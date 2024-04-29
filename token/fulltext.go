package token

func TokenizeFulltext(og []string) []string {
	var toks []string
	for _, v := range og {
		tokens := Split(v)
		tokens = RemoveStopwords(tokens...)
		for _, t := range tokens {
			t = Stem(t)
			toks = append(toks, Normalize(t))
		}
	}
	return toks
}
