package analyzer

import (
	"fmt"
	"testing"
)

const testStr = `Tad's finally entered the ultimate dungeon, Titan—though it wasn't with the fanfare of a triumphant hero, but with the intense desperation of a quiet death.`

func TestSplitOnWhitespace(t *testing.T) {
	splt := splitT()
	if len(splt) != 28 {
		t.Errorf("%v\n", splt)
	}
}

func TestKeywordTokenize(t *testing.T) {
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

	tokens := Keyword.Tokenize(tests...)
	err := testTokenizer(want, tokens)
	if err != nil {
		t.Error(err)
	}
}

func TestFulltextTokenize(t *testing.T) {
	toks := Standard.Tokenize(testStr)
	want := 15
	if len(toks) != want {
		t.Errorf("got %v, wanted %v\n", len(toks), want)
	}
}

func TestSimpleTokenize(t *testing.T) {
	tests := []string{
		"Where All Paths Meet",
		"Mimic & Me",
		"Natural-Born Cullers",
		"The Hitman's Guide to Codenames and Ill-Gotten Gains",
		"HIM",
		"I Knew Him",
		"All in with Him",
		"The Story of the Night",
	}

	want := [][]string{
		[]string{"where", "all", "path", "meet"},
		[]string{"mimic", "me"},
		[]string{"natur", "born", "culler"},
		[]string{"the", "hitman", "guid", "to", "codenam", "and", "ill", "gotten", "gain"},
		[]string{"him"},
		[]string{"i", "knew", "him"},
		[]string{"all", "in", "with", "him"},
		[]string{"the", "stori", "of", "night"},
	}

	for i, test := range tests {
		tokens := Simple.Tokenize(test)
		err := testTokenizer(want[i], tokens)
		if err != nil {
			t.Error(err)
		}
	}
}

func splitT() []string {
	return SplitOnWhitespaceAndPunct(testStr)
}

func testTokenizer(want []string, tokens []string) error {
	for i, tok := range tokens {
		if tok != want[i] {
			return fmt.Errorf("tokens %#v\ngot %v, wanted %s\n", tokens, tok, want[i])
		}
	}
	return nil
}
