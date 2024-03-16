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
	SrchAttr  = `searchableAttributes`
	FacetAttr = `attributesForFaceting`
	SortAttr  = `sortableAttributes`

	// Cfg
	FullText     = `fullText`
	BlvPath      = `fullText`
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
	UID,
}

var paramsCfg = []string{
	DataDir,
	DataFile,
	DefaultField,
	FullText,
	BlvPath,
	UID,
}

var paramsSearch = []string{
	Hits,
	AttributesToRetrieve,
	Page,
	HitsPerPage,
	Query,
	NbHits,
	NbPages,
	SortBy,
	Order,
}

var paramsFacets = []string{
	SortFacetsBy,
	Facets,
	Filters,
	FacetFilters,
	UID,
}
