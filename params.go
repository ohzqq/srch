package srch

const (
	SearchableAttributes  = `searchableAttributes`
	AttributesForFaceting = `attributesForFaceting`
	AttributesToRetrieve  = `attributesToRetrieve`
	Page                  = "page"
	HitsPerPage           = "hitsPerPage"
	SortFacetValuesBy     = `sortFacetValuesBy`
	ParamQuery            = `query`
	ParamFacets           = "facets"
	ParamFacetFilters     = `facetFilters`
	ParamFilters          = "filters"
	DataDir               = `dataDir`
	DataFile              = `dataFile`
	ParamFullText         = `fullText`
)

var ReservedKeys = []string{
	"and",
	"or",
	"field",
	"q",
	"sort_by",
	"order",
	"data_file",
	"data_dir",
	"full_text",
	"query",
	"filters",
	"facetFilters",
}
