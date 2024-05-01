package analyzer

import (
	"strings"
	"unicode"

	"github.com/samber/lo"
)

type Tokenizer struct {
	toLower       bool
	alphaNumOnly  bool
	splitStr      bool
	stem          bool
	rmStopwords   bool
	rmStopLetters bool
	rmPunct       bool
}

func DefaultTokenizer() *Tokenizer {
	return &Tokenizer{
		toLower: true,
	}
}

func (t *Tokenizer) Tokenize(og ...string) []string {
	var tokens []string
	for _, str := range og {
		tokens = append(tokens, t.analyze(str)...)
	}
	return lo.Uniq(tokens)
}

func (t *Tokenizer) analyze(str string) []string {
	str = strings.ToLower(str)

	if !t.splitStr {
		return []string{str}
	}

	tokens := SplitOnWhitespaceAndPunct(str)

	if t.alphaNumOnly {
		tokens = AlphaNumOnly(tokens)
	}

	if t.rmStopwords {
		tokens = RemoveStopWords(tokens...)
	}

	if t.rmStopLetters {
		tokens = RemoveStopLetters(tokens)
	}

	if t.stem {
		tokens = StemWords(tokens)
	}

	return tokens
}

func (t *Tokenizer) split(tok string) []string {
	fn := unicode.IsSpace
	return strings.FieldsFunc(tok, fn)
}

func StemWords(toks []string) []string {
	stem := make([]string, len(toks))
	for i, tok := range toks {
		stem[i] = Stem(tok)
	}
	return stem
}

func ToLower(toks []string) []string {
	low := make([]string, len(toks))
	for i, tok := range toks {
		low[i] = strings.ToLower(tok)
	}
	return low
}

func AlphaNumOnly(toks []string) []string {
	alpha := make([]string, len(toks))
	for i, tok := range toks {
		alpha[i] = AlphaNumericOnly(tok)
	}
	return alpha
}

func RemovePunct(toks []string) []string {
	var none []string
	for _, tok := range toks {
		if len(tok) > 1 {
			none = append(none, tok)
		} else {
			if r := rune(tok[0]); !unicode.IsPunct(r) {
				none = append(none, tok)
			}
		}
	}
	return none
}
