package token

import (
	"testing"
)

const testStr = `Tad's finally entered the ultimate dungeon, Titan—though it was not with the fanfare of a triumphant hero, but with the intense desperation of a quiet death.`

func TestSplitOnWhitespace(t *testing.T) {
	splt := splitT()
	if len(splt) != 28 {
		t.Errorf("%v\n", splt)
	}
}

func TestRemoveStopwords(t *testing.T) {
	w := RemoveStopwords(splitT()...)
	if len(w) != 22 {
		t.Errorf("%v\n", len(w))
	}
}

func TestNormalize(t *testing.T) {
	want := `tad s finally entered the ultimate dungeon titan though it was not with the fanfare of a triumphant hero but with the intense desperation of a quiet death`
	toks := Normalize(testStr)
	if toks != want {
		t.Errorf("got %v, wanted %s\n", toks, want)
	}
}

func TestTokenize(t *testing.T) {
	ft := &Fulltext{}
	toks := ft.Tokenize(testStr)
	want := 22
	if len(toks) != want {
		t.Errorf("got %v, wanted %v\n", len(toks), want)
	}
}

func splitT() []string {
	return Split(testStr)
}
