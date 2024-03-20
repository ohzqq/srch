package param

const (
	// search params
	Hits         = `hits`
	RtrvAttr     = `attributesToRetrieve`
	Page         = "page"
	HitsPerPage  = "hitsPerPage"
	SortFacetsBy = `sortFacetValuesBy`
	Query        = `query`
	Facets       = "facets"
	Filters      = "filters"
	FacetFilters = `facetFilters`
	NbHits       = `nbHits`
	NbPages      = `nbPages`
	SortBy       = `sortBy`
	Order        = `order`

	// Settings
	SrchAttr  = `searchableAttributes`
	FacetAttr = `attributesForFaceting`
	SortAttr  = `sortableAttributes`

	// Cfg
	Format       = `format`
	DefaultField = `title`
	UID          = `uid`

	// file paths
	Route = `route`
	Blv   = "blv"
	Dir   = "dir"
	File  = "file"

	// content-type
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
)

var SettingParams = []string{
	SrchAttr,
	FacetAttr,
	SortAttr,
	UID,
	DefaultField,
	Format,
}

var SearchParams = []string{
	Hits,
	RtrvAttr,
	Page,
	HitsPerPage,
	Query,
	NbHits,
	NbPages,
	SortBy,
	Order,
	SortFacetsBy,
	Facets,
	Filters,
	FacetFilters,
	Route,
}
