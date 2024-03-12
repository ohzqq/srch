package srch

import (
	"fmt"
	"net/url"
	"slices"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/srch/param"
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
	"searchableAttributes=title&fullText=testadata/poot.bleve",
	"searchableAttributes=title&dataDir=testdata/data-dir",
	"attributesForFaceting=tags,authors,series,narrators",
	"attributesForFaceting=tags,authors,series,narrators&dataFile=testdata/data-dir/audiobooks.json",
	"searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators",
	"searchableAttributes=title&dataFile=testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators",
	`searchableAttributes=title&dataFile=testdata/data-dir/audiobooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
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
			"query":   []string{"fish"},
		},
	},
}

func TestNewParser(t *testing.T) {
	tests := []string{
		`searchableAttributes=title&dataDir=testdata/nddata/&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
		`searchableAttributes=title&dataFile=testdata/ndbooks.json&attributesForFaceting=tags,authors,series,narrators&page=3&query=fish&facets=tags&facets=authors&sortBy=title&order=desc&facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
	}
	for _, query := range tests {
		//q := url.QueryEscape(`facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`)
		//query += "&"
		//query += q
		params, err := param.Parse(query)
		if err != nil {
			t.Error(err)
		}
		//println(query)
		data, err := params.Settings.GetData()
		if err != nil {
			t.Error(err)
		}
		if params.HasData() && len(data) != 7252 {
			t.Errorf("got %d, expected %d\n", len(data), 7252)
		}
		//fmt.Printf("Settings %+v\n", params.Settings)
		//fmt.Printf("Search %+v\n", params.Search)
	}
}

func TestMapStruct(t *testing.T) {
	t.SkipNow()
	for _, query := range testQuerySettings {
		vals, err := url.ParseQuery(query)
		if err != nil {
			t.Error(err)
		}

		fmt.Printf("%+v\n", vals)

		params := &param.Params{}

		err = mapstructure.Decode(vals, params)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%+v\n", params)
	}
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
