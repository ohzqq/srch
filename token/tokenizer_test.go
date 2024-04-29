package token

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

func TestRemoveStopwords(t *testing.T) {
	w := RemoveStopwords(splitT()...)
	if len(w) != 15 {
		t.Errorf("%v\n", len(w))
	}
}

func TestNormalize(t *testing.T) {
	want := `tad s finally entered the ultimate dungeon titan though it wasn t with the fanfare of a triumphant hero but with the intense desperation of a quiet death`
	toks := Normalize(testStr)
	if toks != want {
		t.Errorf("got %v\nwanted %s\n", toks, want)
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

	tokens := Keywords.Tokenize(tests...)
	for i, tok := range tokens {
		if tok != want[i] {
			t.Errorf("got %s, wanted %s\n", tok, want[i])
		}
	}
}

func TestFulltextTokenize(t *testing.T) {
	toks := Fulltext.Tokenize(testStr)
	want := 15
	if len(toks) != want {
		t.Errorf("got %v, wanted %v\n", len(toks), want)
	}
}

func TestSimpleTokenize(t *testing.T) {
	toks := Simple.Tokenize(testStr)
	want := 26
	if len(toks) != want {
		fmt.Printf("toks %#v\n", toks)
		t.Errorf("got %v, wanted %v\n", len(toks), want)
	}
}

func splitT() []string {
	return Split(testStr)
}
