package idx

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ohzqq/srch/param"
)

const numBooks = 7251

type paramTest struct {
	query string
	want  *param.Params
}

const (
	facetParamStr   = `facets=tags,authors,series,narrators`
	facetParamSlice = `facets=tags&facets=authors&facets=series&facets=narrators`
	srchAttrParam   = "searchableAttributes=title"
	queryParam      = `query=fish`
	sortParam       = `sortBy=title&order=desc`
	filterParam     = `facetFilters=["authors:amy lane", ["tags:romance", "tags:-dnr"]]`
	uidParam        = `uid=id`
)

const (
	testDataFile   = `testdata/nddata/ndbooks.ndjson`
	testDataDir    = `testdata/data-dir`
	testBlvPath    = `testdata/poot.bleve`
	testHareDskDir = `/home/mxb/code/srch/param/testdata/hare`
	testHareURL    = `/testdata/hare`
)

var testQuerySettings = []string{
	``,
	mkURL("", `searchableAttributes`),
	mkURL("", srchAttrParam),
	dirRoute(srchAttrParam),
	mkURL("", facetParamStr),
	dirRoute(facetParamSlice),
	mkURL("", srchAttrParam, facetParamStr),
	dirRoute(srchAttrParam, facetParamStr),
}

var paramTests = []paramTest{
	paramTest{
		query: testQuerySettings[0],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: testQuerySettings[1],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: testQuerySettings[2],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: testQuerySettings[3],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{""},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: testQuerySettings[4],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{"tags", "authors", "series", "narrators"},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: testQuerySettings[5],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"*"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{"tags", "authors", "series", "narrators"},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: testQuerySettings[6],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{"tags", "authors", "series", "narrators"},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: testQuerySettings[7],
		want: &param.Params{
			Other:        url.Values{},
			SrchAttr:     []string{"title"},
			FacetAttr:    []string{"tags", "authors", "series", "narrators"},
			SortAttr:     []string{""},
			DefaultField: "",
			Hits:         0,
			RtrvAttr:     []string{""},
			Page:         0,
			HitsPerPage:  0,
			Query:        "",
			SortBy:       "",
			Order:        "",
			Facets:       []string{"tags", "authors", "series", "narrators"},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
}

func checkIdxName(got, want string) error {
	if got != want {
		return fmt.Errorf("index name is %v, wanted %v\n", got, want)
	}
	return nil
}

func checkAttrs(field string, attrs []string, want []string) error {
	if !slices.Equal(attrs, want) {
		return fmt.Errorf("for %v got %#v, wanted %#v\n", field, attrs, want)
	}
	return nil
}

func dirRoute(params ...string) string {
	params = append(params, "path="+testDataDir)
	return mkURL(param.Dir.String(), params...)
}

func fileRoute(params ...string) string {
	params = append(params, "path="+testDataFile)
	return mkURL(param.File.String(), params...)
}

func hareURL(params ...string) string {
	params = append(params, "path="+mkHarePath("index"))
	return mkURL(param.File.String(), params...)
}

func totalBooksErr(total int, vals ...any) error {
	if total != numBooks && total != 7251 {
		err := fmt.Errorf("got %d, expected %d\n", total, numBooks)
		return fmt.Errorf("%w\nmsg: %v", err, vals)
	}
	return nil
}

func mkURL(path string, rq ...string) string {
	u := &url.URL{
		Path:     path,
		RawQuery: strings.Join(rq, "&"),
	}
	return "?" + u.RawQuery
}

func mkHarePath(name string) string {
	return filepath.Join(testHareURL, name+".json")
}
