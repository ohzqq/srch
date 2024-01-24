package field

import "testing"

func TestStripAlphaNum(t *testing.T) {
	t.SkipNow()
	names := []any{
		"holiday.christmas",
		"grumpy/sunshine",
		"L.A. Witt",
		"Breath & Fire",
		"[Psychokinetic] Eyeball Pulling",
	}

	tokens := KeywordAnalyzer(names)
	for _, token := range tokens {
		println(token)
	}
}
