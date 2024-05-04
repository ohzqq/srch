package db

import (
	"slices"
	"testing"

	"github.com/ohzqq/srch/analyzer"
)

var testMapping = map[analyzer.Analyzer][]string{
	analyzer.Standard: []string{"comments", "tags"},
	analyzer.Simple:   []string{"title"},
	analyzer.Keyword:  []string{"tags", "authors", "series", "narrators"},
}

func TestNewMapping(t *testing.T) {
	m := NewMapping()
	m.AddFulltext("comments", "tags")
	m.AddSimple("title")
	m.AddKeywords("tags", "authors", "series", "narrators")
	for ana, fields := range m {
		want := testMapping[ana]
		if !slices.Equal(fields, want) {
			t.Errorf("got %#v\n, wanted %#v\n", fields, want)
		}
	}
}

func TestMappingParams(t *testing.T) {
	var testMapping = map[analyzer.Analyzer][]string{
		analyzer.Standard: []string{"title", "comments", "tags"},
		analyzer.Keyword:  []string{"tags", "authors", "series", "narrators"},
	}

	m := NewMappingFromParams(testParams())

	for ana, fields := range m {
		want := testMapping[ana]
		if !slices.Equal(fields, want) {
			t.Errorf("got %#v\n, wanted %#v\n", fields, want)
		}
	}
}
