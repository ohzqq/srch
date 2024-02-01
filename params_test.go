package srch

import (
	"fmt"
	"net/url"
	"slices"
	"testing"
)

var (
	defFields    = []string{DefaultField}
	defDataDir   = []string{testDataDir}
	defDataFile  = []string{testDataFile}
	defFacetAttr = []string{
		"tags",
		"authors",
		"series",
		"narrators",
	}
)

var testQuerySettings = []string{
	"",
	"searchableAttributes=",
	"searchableAttributes=title&fullText",
	"searchableAttributes=title&dataDir=testdata/data-dir",
	"attributesForFaceting=tags,authors,series,narrators",
	"attributesForFaceting=tags,authors,series,narrators&dataFile=testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators",
	"searchableAttributes=title&dataFile=testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators",
}

var testParsedParams = []*Params{
	&Params{
		Settings: url.Values{
			SrchAttr: defFields,
		},
	},
	&Params{
		Settings: url.Values{
			SrchAttr: defFields,
		},
	},
	&Params{
		Settings: url.Values{
			SrchAttr: defFields,
		},
	},
	&Params{
		Settings: url.Values{
			SrchAttr: defFields,
			DataDir:  []string{testDataDir},
		},
	},
	&Params{
		Settings: url.Values{
			SrchAttr:  defFields,
			FacetAttr: defFacetAttr,
		},
	},
	&Params{
		Settings: url.Values{
			SrchAttr:  defFields,
			DataFile:  defDataFile,
			FacetAttr: defFacetAttr,
		},
	},
	&Params{
		Settings: url.Values{
			SrchAttr:  defFields,
			FacetAttr: defFacetAttr,
		},
	},
	&Params{
		Settings: url.Values{
			SrchAttr:  defFields,
			DataFile:  defDataFile,
			FacetAttr: defFacetAttr,
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

	if got.Settings.Has(DataDir) {
		vals := got.GetSlice(DataDir)
		if !slices.Equal(vals, want.GetSlice(DataDir)) {
			err = fmt.Errorf(fmtStr, err, vals, want.GetSlice(DataDir))
		}
	}

	if got.Settings.Has(DataFile) {
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
