package param

import (
	"strings"

	"github.com/gobuffalo/flect"
)

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

// Query returns a camelCase string to use as the key in a URL query
func (p Param) Query() string {
	return p.String()
}

// Snake returns an alphanumeric, lowercased, underscored string
func (p Param) Snake() string {
	return flect.Underscore(p.String())
}

// ToLower returns an all lowercase string
func (p Param) ToLower() string {
	return strings.ToLower(p.String())
}

// Dasherize returns an alphanumeric, lowercased, dashed string
func (p Param) Dasherize() string {
	return flect.Dasherize(p.String())
}

// Slug returns an alphanumeric, lowercased, dashed string
func (p Param) Slug() string {
	return p.Dasherize()
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
