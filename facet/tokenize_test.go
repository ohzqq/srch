package facet

import (
	"fmt"
	"testing"
)

const testStr = `Tad's finally entered the ultimate dungeon, Titan—though it was not with the fanfare of a triumphant hero, but with the intense desperation of a quiet death.`

func TestFieldsFunc(t *testing.T) {
	splt := splitT()
	if len(splt) != 28 {
		t.Errorf("%v\n", splt)
	}
}

func TestRemoveStopwords(t *testing.T) {
	w := removeStopwords(splitT())
	if len(w) != 22 {
		t.Errorf("%v\n", len(w))
	}
}

func TestNormalize(t *testing.T) {
	w := removeStopwords(splitT())
	for _, tok := range w {
		norm := normalizeStr(tok)
		println(norm)
	}
}

func TestTokenize(t *testing.T) {
	toks := Tokenize(testStr)
	println(len(toks))
	fmt.Printf("%v\n", toks)
}

func splitT() []string {
	return split(testStr)
}
