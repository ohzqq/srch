package analyzer

import (
	"fmt"
	"testing"
)

func TestTokenizerKeyword(t *testing.T) {
	tests := []string{
		"grumpy/sunshine",
		"best friend's brother",
		"angst",
		"ABO",
	}
	want := []string{
		"grumpy/sunshine",
		"best friend's brother",
		"angst",
		"abo",
	}
	tokenizer := Keyword.tokenizer()
	tokens := tokenizer.Tokenize(tests...)
	err := testTokenizer(want, tokens)
	if err != nil {
		t.Error(err)
	}
}

func TestTokenizerSimple(t *testing.T) {
	tests := []string{
		"Where All Paths Meet",
		"Mimic & Me",
		"Natural-Born Cullers",
		"The Hitman's Guide to Codenames and Ill-Gotten Gains",
		"HIM",
		"I Knew Him",
		"All in with Him",
	}

	want := [][]string{
		[]string{"where", "all", "path", "meet"},
		[]string{"mimic", "me"},
		[]string{"natur", "born", "culler"},
		[]string{"the", "hitman", "guid", "to", "codenam", "and", "ill", "gotten", "gain"},
		[]string{"him"},
		[]string{"i", "knew", "him"},
		[]string{"all", "in", "with", "him"},
	}

	tokenizer := Simple.tokenizer()
	for i, test := range tests {
		tokens := tokenizer.Tokenize(test)
		err := testTokenizer(want[i], tokens)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestTokenizerStandard(t *testing.T) {
}

func testTokenizer(want []string, tokens []string) error {
	for i, tok := range tokens {
		if tok != want[i] {
			return fmt.Errorf("got %v, wanted %s\n", tok, want[i])
		}
	}
	return nil
}
