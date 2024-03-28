package param

import "github.com/gobuffalo/flect"

type Param int

//go:generate stringer -type=Param -linecomment
const (
	// search params
	Hits         Param = iota // hits
	RtrvAttr                  // attributesToRetrieve
	Page                      // page
	HitsPerPage               // hitsPerPage
	SortFacetsBy              // sortFacetValuesBy
	MaxFacetVals              // maxValuesPerFacet
	Query                     // query
	Facets                    // facets
	Filters                   // filters
	FacetFilters              // facetFilters
	NbHits                    // nbHits
	NbPages                   // nbPage
	SortBy                    // sortBy
	Order                     // order

	// Settings
	SrchAttr  // searchableAttributes
	FacetAttr // attributesForFaceting
	SortAttr  // sortableAttributes
	Path      // path

	// Cfg
	Format       // format
	DefaultField // title
	UID          // uid

	// file paths
	Route // route
	Blv   // blv
	Dir   // dir
	File  // file
)

const (
	// content-type
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
)

func (p Param) Query() string {
	return p.String()
}

func (p Param) Snake() string {
	return flect.Underscore(p.String())
}

var SettingParams = []Param{
	SrchAttr,
	FacetAttr,
	SortAttr,
	UID,
	DefaultField,
	Format,
}

var SearchParams = []Param{
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
	MaxFacetVals,
	Facets,
	Filters,
	FacetFilters,
	Route,
	Path,
}

var Routes = []Param{
	Blv,
	Dir,
	File,
}
