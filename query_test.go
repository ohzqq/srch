package srch

import (
	"net/url"
	"slices"
	"testing"
)

const testValuesCfg = `and=tags:count:desc&field=title&or=authors:label:asc&or=narrators&or=series&data_file=testdata/data-dir/audiobooks.json&sort_by=title`

var testQueryStrings = []string{
	"",
	"searchableAttributes=",
	"searchableAttributes=title",
	"searchableAttributes=title&dataDir=testdata/data-dir",
	"attributesForFaceting=tags,authors,series,narrators",
	"attributesForFaceting=tags,authors,series,narrators&dataFile=testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators",
	"searchableAttributes=title&dataFile=testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators",
}

var testParsedParams = []*Params{
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
		},
	},
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
		},
	},
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
		},
	},
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
			DataDir:  []string{testDataDir},
		},
	},
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
			FacetAttr: []string{
				"tags",
				"authors",
				"series",
				"narrators",
			},
		},
	},
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
			DataFile: []string{testData},
			FacetAttr: []string{
				"tags",
				"authors",
				"series",
				"narrators",
			},
		},
	},
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
			FacetAttr: []string{
				"tags",
				"authors",
				"series",
				"narrators",
			},
		},
	},
	&Params{
		Values: url.Values{
			SrchAttr: []string{DefaultField},
			DataFile: []string{testData},
			FacetAttr: []string{
				"tags",
				"authors",
				"series",
				"narrators",
			},
		},
	},
}

func TestParseQueryStrings(t *testing.T) {
	for i, q := range testQueryStrings {
		p := NewQuery(q)
		want := testParsedParams[i]
		if attr := p.SrchAttr(); !slices.Equal(attr, want.SrchAttr()) {
			t.Errorf("query: %s\ngot %v, expected %v\n", q, attr, want.SrchAttr())
		}
	}
}

func queryParamsValues() url.Values {
	vals := make(url.Values)
	vals.Set("facetFilters", testComboFilter())
	return vals
}

func queryParamsString() string {
	return queryParamsValues().Encode()
}

func requestParams() string {
	p := queryParamsString()
	return p
}

func getNewQuery() url.Values {
	return ParseQuery(testValuesCfg, testQueryString, testSearchString)
}
