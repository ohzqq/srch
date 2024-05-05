package idx

import (
	"fmt"
	"net/url"
	"slices"
)

type paramTest struct {
	query string
	want  *Params
}

var paramTests = []paramTest{
	paramTest{
		query: ``,
		want: &Params{
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
			Index:        "index",
		},
	},
	paramTest{
		query: `?searchableAttributes=`,
		want: &Params{
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
			Index:        "index",
		},
	},
	paramTest{
		query: `?searchableAttributes=title`,
		want: &Params{
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
			Index:        "index",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&dataDir=../testdata/hare`,
		want: &Params{
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
			Index:        "index",
		},
	},
	paramTest{
		query: `?attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
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
			Index:        "index",
		},
	},
	paramTest{
		query: `?attributesForFaceting=tags,authors,series,narrators&dataDir=../testdata/hare`,
		want: &Params{
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
			Index:        "index",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators`,
		want: &Params{
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
			Index:        "index",
		},
	},
	paramTest{
		query: `?searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&dataDir=../testdata/hare`,
		want: &Params{
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
			Index:        "index",
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
		return fmt.Errorf("for %v got %#v, wanted %#v\n", field, attrs, idx)
	}
	return nil
}
