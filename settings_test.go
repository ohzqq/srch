package srch

import (
	"slices"
	"testing"
)

type settingsTest struct {
	query string
	want  *Settings
}

var settingsTestVals = []settingsTest{
	settingsTest{
		query: "",
		want: &Settings{
			SearchableAttributes: []string{"title"},
		},
	},
	settingsTest{
		query: "searchableAttributes=",
		want: &Settings{
			SearchableAttributes: []string{"title"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title&data_file=testdata/data-dir/audiobooks.json",
		want: &Settings{
			SearchableAttributes: []string{"title"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title,comments",
		want: &Settings{
			SearchableAttributes: []string{"title", "comments"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title&searchableAttributes=comments",
		want: &Settings{
			SearchableAttributes: []string{"title", "comments"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title&attributesForFaceting=",
		want: &Settings{
			SearchableAttributes: []string{"title"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title&attributesForFaceting=tags&data_file=testdata/data-dir/audiobooks.json",
		want: &Settings{
			SearchableAttributes:  []string{"title"},
			AttributesForFaceting: []string{"tags"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title&attributesForFaceting=tags,authors&data_file=testdata/data-dir/audiobooks.json",
		want: &Settings{
			SearchableAttributes:  []string{"title"},
			AttributesForFaceting: []string{"tags", "authors"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title&attributesForFaceting=tags&attributesForFaceting=authors",
		want: &Settings{
			SearchableAttributes:  []string{"title"},
			AttributesForFaceting: []string{"tags", "authors"},
		},
	},
	settingsTest{
		query: "searchableAttributes=title&full_text=",
		want: &Settings{
			SearchableAttributes: []string{"title"},
			TextAnalyzer:         Text,
		},
	},
	settingsTest{
		query: "searchableAttributes=title&attributesForFaceting=tags&full_text",
		want: &Settings{
			SearchableAttributes:  []string{"title"},
			AttributesForFaceting: []string{"tags"},
			TextAnalyzer:          Text,
		},
	},
}

func TestSettingsValueParser(t *testing.T) {
	for _, test := range settingsTestVals {
		q := ParseQuery(test.query)

		sea := GetQueryStringSlice("searchableAttributes", q)
		if !slices.Equal(sea, test.want.SearchableAttributes) {
			t.Errorf("%s: got %+v, wanted %+v\n", test.query, sea, test.want.SearchableAttributes)
		}
		facet := GetQueryStringSlice("attributesForFaceting", q)
		if !slices.Equal(facet, test.want.AttributesForFaceting) {
			t.Errorf("%s: got %+v, wanted %+v\n", test.query, facet, test.want.AttributesForFaceting)
		}
		analyzer := GetAnalyzer(q)
		if analyzer != Text && test.want.TextAnalyzer == Text {
			t.Errorf("%s: got %+v, wanted %+v\n", test.query, analyzer, test.want.TextAnalyzer)
		}
	}
}

func TestSettings(t *testing.T) {
	for _, test := range settingsTestVals {
		settings := NewSettings(test.query)
		if !slices.Equal(settings.SearchableAttributes, test.want.SearchableAttributes) {
			t.Errorf("%s: got %+v, wanted %+v\n", test.query, settings.SearchableAttributes, test.want.SearchableAttributes)
		}
		if !slices.Equal(settings.AttributesForFaceting, test.want.AttributesForFaceting) {
			t.Errorf("%s: got %+v, wanted %+v\n", test.query, settings.AttributesForFaceting, test.want.AttributesForFaceting)
		}
		if settings.TextAnalyzer != Text && test.want.TextAnalyzer == Text {
			t.Errorf("%s: got %+v, wanted %+v\n", test.query, settings.TextAnalyzer, test.want.TextAnalyzer)
		}
	}
}
