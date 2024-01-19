package srch

import (
	"slices"
	"testing"
)

var settingsTestVals = map[string]*Settings{
	"": &Settings{
		SearchableAttributes: []string{"title"},
	},
	"searchableAttributes=": &Settings{
		SearchableAttributes: []string{"title"},
	},
	"searchableAttributes=title": &Settings{
		SearchableAttributes: []string{"title"},
	},
	"searchableAttributes=title,comments": &Settings{
		SearchableAttributes: []string{"title", "comments"},
	},
	"searchableAttributes=title&searchableAttributes=comments": &Settings{
		SearchableAttributes: []string{"title", "comments"},
	},
	"searchableAttributes=title&attributesForFaceting=": &Settings{
		SearchableAttributes: []string{"title"},
	},
	"searchableAttributes=title&attributesForFaceting=tags": &Settings{
		SearchableAttributes:  []string{"title"},
		AttributesForFaceting: []string{"tags"},
	},
	"searchableAttributes=title&attributesForFaceting=tags,series": &Settings{
		SearchableAttributes:  []string{"title"},
		AttributesForFaceting: []string{"tags", "series"},
	},
	"searchableAttributes=title&attributesForFaceting=tags&attributesForFaceting=series": &Settings{
		SearchableAttributes:  []string{"title"},
		AttributesForFaceting: []string{"tags", "series"},
	},
}

func TestSettings(t *testing.T) {
	for test, want := range settingsTestVals {
		settings := NewSettings(test)
		if !slices.Equal(settings.SearchableAttributes, want.SearchableAttributes) {
			t.Errorf("%s: got %+v, wanted %+v\n", test, settings.SearchableAttributes, want.SearchableAttributes)
		}
		if !slices.Equal(settings.AttributesForFaceting, want.AttributesForFaceting) {
			t.Errorf("%s: got %T, wanted %T\n", test, settings.AttributesForFaceting, want.AttributesForFaceting)
		}
	}
}
