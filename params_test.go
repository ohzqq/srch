package srch

import (
	"fmt"
	"net/url"
	"slices"
	"testing"
)

const testValuesCfg = `and=tags:count:desc&field=title&or=authors:label:asc&or=narrators&or=series&data_file=testdata/data-dir/audiobooks.json&sort_by=title`

var testQuerySettings = []string{
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
	for i, q := range testQuerySettings {
		p := ParseParams(q)
		want := testParsedParams[i]
		if err := settingsErr(p, want); err != nil {
			t.Error(err)
		}
	}
}

func settingsErr(got *Params, want *Params) error {
	var err error
	fmtStr := "%w\ngot %v, expected %v\n"
	attr := got.SrchAttr()
	if !slices.Equal(attr, want.SrchAttr()) {
		err = fmt.Errorf(fmtStr, err, attr, want.SrchAttr())
	}

	facet := got.FacetAttr()
	if !slices.Equal(facet, want.FacetAttr()) {
		err = fmt.Errorf(fmtStr, err, facet, want.FacetAttr())
	}

	if got.Values.Has(DataDir) {
		vals := got.GetSlice(DataDir)
		if !slices.Equal(vals, want.GetSlice(DataDir)) {
			err = fmt.Errorf(fmtStr, err, vals, want.GetSlice(DataDir))
		}
	}

	if got.Values.Has(DataFile) {
		vals := got.GetSlice(DataFile)
		if !slices.Equal(vals, want.GetSlice(DataFile)) {
			err = fmt.Errorf(fmtStr, err, vals, want.GetSlice(DataFile))
		}
	}
	return err
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