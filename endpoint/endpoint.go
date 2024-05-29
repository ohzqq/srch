package endpoint

import (
	"path/filepath"
)

const (
	IdxName   = `indexName`
	ObjectID  = `objectID`
	FacetName = `facetName`
)

type endpoint int

const (
	Root endpoint = iota
	Idx
	IdxBatch
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
	routeIdxBatch    = "/indexes/{indexName}/batch"
	routeIdxBrowse   = "/indexes/{indexName}/browse"
	routeIdxObject   = "/indexes/{indexName}/{objectID}"
	routeIdxQuery    = "/indexes/{indexName}/query"
	routeIdxSettings = "/indexes/{indexName}/settings"
	routeFacets      = "/indexes/{indexName}/facets"
	routeFacet       = "/indexes/{indexName}/facets/{facetName}"
	routeFacetQuery  = "/indexes/{indexName}/facets/{facetName}/query"
)

var Endpoints = []endpoint{
	Root,
	Idx,
	IdxBatch,
	IdxBrowse,
	IdxObject,
	IdxQuery,
	IdxSettings,
	Facets,
	Facet,
	FacetQuery,
}

func (end endpoint) SetWildcards(sets ...string) string {
	if end != Root {
		u := Root.Route()
		if len(sets) > 0 {
			u = filepath.Join(u, sets[0])
			switch end {
			case Idx:
				return u
			case IdxBatch:
				return filepath.Join(u, "batch")
			case IdxBrowse:
				return filepath.Join(u, "browse")
			case IdxQuery:
				return filepath.Join(u, "query")
			case IdxSettings:
				return filepath.Join(u, "settings")
			case Facets:
				return filepath.Join(u, "facets")
			}
		}
		if len(sets) > 1 {
			switch end {
			case Facet:
				return filepath.Join(u, "facets", sets[1])
			case FacetQuery:
				return filepath.Join(u, "facets", sets[1], "query")
			case IdxObject:
				return filepath.Join(u, sets[1])
			}
		}
	}
	return Root.Route()
}

func (end endpoint) Route() string {
	switch end {
	case Idx:
		return routeIdx
	case IdxBatch:
		return routeIdxBatch
	case IdxBrowse:
		return routeIdxBrowse
	case IdxObject:
		return routeIdxObject
	case IdxQuery:
		return routeIdxQuery
	case IdxSettings:
		return routeIdxSettings
	case Facets:
		return routeFacets
	case Facet:
		return routeFacet
	case FacetQuery:
		return routeFacetQuery
	default:
		return routeIndexes
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

var Routes = []string{
	routeIndexes,
	routeIdx,
	routeIdxBrowse,
	routeIdxObject,
	routeIdxQuery,
	routeIdxSettings,
	routeFacets,
	routeFacet,
	routeFacetQuery,
}
