package param

const (
	// search params
	Hits                 = `hits`
	AttributesToRetrieve = `attributesToRetrieve`
	Page                 = "page"
	HitsPerPage          = "hitsPerPage"
	SortFacetsBy         = `sortFacetValuesBy`
	Query                = `query`
	Facets               = "facets"
	Filters              = "filters"
	FacetFilters         = `facetFilters`
	NbHits               = `nbHits`
	NbPages              = `nbPages`
	SortBy               = `sortBy`
	Order                = `order`

	// Settings
	FullText     = `fullText`
	SrchAttr     = `searchableAttributes`
	FacetAttr    = `attributesForFaceting`
	SortAttr     = `sortableAttributes`
	DataDir      = `dataDir`
	DataFile     = `dataFile`
	DefaultField = `title`
	UID          = `uid`

	TextAnalyzer    = "text"
	KeywordAnalyzer = "keyword"
)

var paramsSettings = []string{
	SrchAttr,
	FacetAttr,
	SortAttr,
	DataDir,
	DataFile,
	DefaultField,
	FullText,
	UID,
}

var paramsSearch = []string{
	Hits,
	AttributesToRetrieve,
	Page,
	HitsPerPage,
	SortFacetsBy,
	Query,
	Facets,
	Filters,
	FacetFilters,
	NbHits,
	NbPages,
	SortBy,
	Order,
}
