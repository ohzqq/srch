// Code generated by "stringer -type=Param -linecomment"; DO NOT EDIT.

package param

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Hits-0]
	_ = x[RtrvAttr-1]
	_ = x[Page-2]
	_ = x[HitsPerPage-3]
	_ = x[SortFacetsBy-4]
	_ = x[MaxFacetVals-5]
	_ = x[Query-6]
	_ = x[Facets-7]
	_ = x[Filters-8]
	_ = x[FacetFilters-9]
	_ = x[NbHits-10]
	_ = x[NbPages-11]
	_ = x[SortBy-12]
	_ = x[Order-13]
	_ = x[SrchAttr-14]
	_ = x[FacetAttr-15]
	_ = x[SortAttr-16]
	_ = x[Path-17]
	_ = x[Format-18]
	_ = x[DefaultField-19]
	_ = x[UID-20]
	_ = x[Route-21]
	_ = x[Blv-22]
	_ = x[Dir-23]
	_ = x[File-24]
	_ = x[Index-25]
}

const _Param_name = "hitsattributesToRetrievepagehitsPerPagesortFacetValuesBymaxValuesPerFacetqueryfacetsfiltersfacetFiltersnbHitsnbPagesortByordersearchableAttributesattributesForFacetingsortableAttributespathformattitleuidrouteblvdirfileindex"

var _Param_index = [...]uint8{0, 4, 24, 28, 39, 56, 73, 78, 84, 91, 103, 109, 115, 121, 126, 146, 167, 185, 189, 195, 200, 203, 208, 211, 214, 218, 223}

func (i Param) String() string {
	if i < 0 || i >= Param(len(_Param_index)-1) {
		return "Param(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Param_name[_Param_index[i]:_Param_index[i+1]]
}
