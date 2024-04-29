package token

type Keyword struct {
	og string
}

func (k *Keyword) Tokenize(og string) []string {
	k.og = og
}
