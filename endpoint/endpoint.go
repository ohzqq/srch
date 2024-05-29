package endpoint

import "path/filepath"

const (
	IdxName   = `{indexName}`
	ObjectID  = `{objectID}`
	FacetName = `{facetName}`
)

type endpoint int

const (
	Indexes endpoint = iota
	Idx
	IdxBrowse
	IdxObject
	IdxQuery
	IdxSettings
	Facets
	Facet
	FacetQuery
)

const (
	routeIndexes     = "/indexes"
	routeIdx         = "/indexes/{indexName}"
	routeIdxBrowse   = "/indexes/{indexName}/browse"
	routeIdxObject   = "/indexes/{indexName}/{objectID}"
	routeIdxQuery    = "/indexes/{indexName}/query"
	routeIdxSettings = "/indexes/{indexName}/settings"
	routeFacets      = "/indexes/{indexName}/facets"
	routeFacet       = "/indexes/{indexName}/facets/{facetName}"
	routeFacetQuery  = "/indexes/{indexName}/facets/{facetName}/query"
)

func (end endpoint) SetWildcards(sets ...string) string {
	if end != Indexes {
		u := Indexes.Route()
		if len(sets) > 0 {
			u = filepath.Join(u, sets[0])
		}
		if len(sets) > 1 {
			u = filepath.Join(u, sets[1])
		}
		return u
	}
	return Indexes.Route()
}

func (end endpoint) Route() string {
	switch end {
	case Indexes:
		return routeIndexes
	case Idx:
		return routeIdx
	case IdxBrowse:
		return routeIdxBrowse
	case IdxObject:
		return routeIdxObject
	case IdxQuery:
		return routeIdxQuery
	case IdxSettings:
		return routeIdxSettings
	case IdxFacets:
		return routeFacets
	case IdxFacet:
		return routeFacet
	case IdxFacetQuery:
		return routeFacetQuery
	}
}

func (end endpoint) Get() string {
	return "GET " + end.Route()
}

func (end endpoint) Put() string {
	return "PUT " + end.Route()
}

func (end endpoint) Del() string {
	return "DELETE " + end.Route()
}

func (end endpoint) Post() string {
	return "POST " + end.Route()
}
