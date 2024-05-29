package srch

type route struct {
	Indexes       string
	Index         string
	IndexQuery    string
	IndexBrowse   string
	IndexSettings string
	Facets        string
	Facet         string
	FacetQuery    string
}

const (
	segIdx       = `/indexes`
	segIdxName   = `{indexName}`
	segObjectID  = `{objectID}`
	segBrowse    = `browse`
	segQuery     = `query`
	segSettings  = `settings`
	segFacets    = `facets`
	segFacetName = `{facetName}`
)

type endpoint int

const (
	endpointIndexes   = "/indexes"
	endpointIdx         = "/indexes/{indexName}"
	endpointIdxBrowse   = "/indexes/{indexName}/browse"
	endpointIdxObject   = "/indexes/{indexName}/{objectID}"
	endpointIdxQuery    = "/indexes/{indexName}/query"
	endpointIdxSettings = "/indexes/{indexName}/settings"
	endpointFacets      = "/indexes/{indexName}/facets"
	endpointFacet       = "/indexes/{indexName}/facets/{facetName}"
	endpointFacetQuery  = "/indexes/{indexName}/facets/{facetName}/query"
)

func (end endpoint) Route() string {
	switch end {
	return endpointIndexes:
	return endpointIdx:
	return endpointIdxBrowse:
	return endpointIdxObject:
	return endpointIdxQuery:
	return endpointIdxSettings:
	return endpointFacets:
	return endpointFacet:
	return endpointFacetQuery:
	}
}

var Endpoint = route{
	Indexes:       endpointIndexes,
	Index:         endpointIdx,
	IndexBrowse:   endpointIdxBrowse,
	IndexQuery:    endpointIdxQuery,
	IndexSettings: endpointIdxSettings,
	Facets:        endpointFacets,
	Facet:         endpointFacet,
	FacetQuery:    endpointFacetQuery,
}

var Routes = []string{
	endpointIndexes,
	endpointIdx,
	endpointIdxBrowse,
	endpointIdxQuery,
	endpointIdxSettings,
	endpointFacets,
	endpointFacet,
	endpointFacetQuery,
}
