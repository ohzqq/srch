package idx

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/ohzqq/srch/param"
)

type paramTest struct {
	query string
	want  *param.Params
}

var paramTests = []paramTest{
	paramTest{
		query: ``,
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
		query: `?searchableAttributes=`,
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
		query: `?searchableAttributes=title`,
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
		query: `?searchableAttributes=title&dataDir=testdata/hare`,
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
		query: `?attributesForFaceting=tags,authors,series,narrators`,
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
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: `?attributesForFaceting=tags,authors,series,narrators&dataDir=testdata/hare`,
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
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
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
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&dataDir=testdata/hare`,
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
			Facets:       []string{""},
			Filters:      "",
			FacetFilters: []any{""},
			SortFacetsBy: "",
			UID:          "",
			IndexName:    "index",
		},
	},
}

func checkIdxName(idx *Idx, want string) error {
	if idx.Name != want {
		return fmt.Errorf("index name is %v, wanted %v\n", idx.Name, want)
	}
	return nil
}

func checkAttrs(field string, attrs []string, want []string) error {
	if !slices.Equal(attrs, want) {
		return fmt.Errorf("for %v got %#v, wanted %#v\n", field, attrs, want)
	}
	return nil
}
