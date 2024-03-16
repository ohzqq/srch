package srch

import (
	"testing"

	"github.com/ohzqq/srch/param"
)

var testFields = []string{
	param.SrchAttr,
	param.FullText,
	param.DataFile,
	param.DataDir,
	param.Facets,
	param.Page,
	param.Query,
	param.SortBy,
	param.Order,
	param.FacetFilters,
}

var reqTests = []map[string]bool{
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     false,
		param.DataDir:      false,
		param.Facets:       false,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     false,
		param.DataDir:      false,
		param.Facets:       false,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     true,
		param.DataFile:     false,
		param.DataDir:      false,
		param.Facets:       false,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     false,
		param.DataDir:      true,
		param.Facets:       false,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     false,
		param.DataDir:      false,
		param.Facets:       true,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     true,
		param.DataDir:      false,
		param.Facets:       true,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     false,
		param.DataDir:      false,
		param.Facets:       true,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     true,
		param.DataDir:      false,
		param.Facets:       true,
		param.Page:         false,
		param.Query:        false,
		param.SortBy:       false,
		param.Order:        false,
		param.FacetFilters: false,
	},
	map[string]bool{
		param.SrchAttr:     true,
		param.FullText:     false,
		param.DataFile:     true,
		param.DataDir:      false,
		param.Facets:       true,
		param.Page:         true,
		param.Query:        true,
		param.SortBy:       true,
		param.Order:        true,
		param.FacetFilters: true,
	},
}

func TestNewRequest(t *testing.T) {
	//for i, test := range testQuerySettings {
	//  req, err := param.Parse(test)
	//  if err != nil {
	//    t.Error(err)
	//  }
	//  w := reqTests[i]
	//  for k, ok := range w {
	//    if ok != req.Has(k) {
	//      t.Errorf("%s:\n %s: got %v, expected %v\n", test, k, req.Has(k), ok)
	//    }
	//  }
	//}

}
