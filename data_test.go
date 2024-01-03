package srch

import "testing"

func TestNewSrc(t *testing.T) {
	testNewSrc(t)
}

func TestFuzzyFindSrc(t *testing.T) {
	src := testNewSrc(t)
	m, err := src.Search("fish")
	if err != nil {
		t.Error(err)
	}
	if len(m.Data) != 56 {
		t.Errorf("got %d, expected 56\n", len(m.Data))
	}

	m, err = src.Search("")
	if err != nil {
		t.Error(err)
	}
	if len(m.Data) != 7174 {
		t.Errorf("got %d, expected 7174\n", len(m.Data))
	}
}

func testNewSrc(t *testing.T) *Src {
	src := NewSrc(books)
	//if src.Len() != 7174 {
	//  t.Errorf("got %d, expected 7174\n", src.Len())
	//}

	return src
}
