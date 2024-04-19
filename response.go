package srch

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/go-http-utils/headers"
	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Response struct {
	*param.Params

	results []map[string]any

	RawQuery    string         `json:"params"`
	FacetFields []*facet.Facet `json:"facetFields"`
	//Facets      *facet.Fields    `json:"facets"`
	facets  *facet.Fields             `json:"facets"`
	Facets  map[string]map[string]int `json:"facets"`
	Hits    []map[string]any          `json:"hits"`
	NbHits  int                       `json:"nbHits"`
	NbPages int                       `json:"nbPages"`
}

func NewResponse(hits []map[string]any, params *param.Params) (*Response, error) {
	res := &Response{
		results:  hits,
		Params:   params,
		RawQuery: params.Encode(),
		Facets:   make(map[string]map[string]int),
	}

	if len(hits) == 0 {
		return res, nil
	}

	if params.Has(param.Facets) {
		facets, err := facet.New(hits, params)
		if err != nil {
			return nil, fmt.Errorf("response failed to calculate facets: %w\n", err)
		}
		res.FacetFields = facets.Facets
		res.facets = facets
		res.results = res.FilterResults()

		for _, facet := range res.FacetFields {
			facet.Items = facet.Keywords()
			facet.Count = facet.Len()
			items := make(map[string]int)
			for _, item := range facet.Items {
				items[item.Label] = item.Count
			}
			res.Facets[facet.Attribute] = items
		}
	}

	res.HitsPerPage = res.hitsPerPage()
	res.NbHits = res.nbHits()
	res.calculatePagination()

	if res.SortBy != "" {
		for _, attr := range res.SortAttr {
			by := NewSort(attr)
			if res.SortBy == by.Field {
				res.results = by.Sort(res.results)
				if res.Order == "desc" {
					slices.Reverse(res.results)
				}
			}
		}
	}

	res.Hits = res.visibleHits()

	return res, nil
}

func (res *Response) calculatePagination() *Response {
	res.HitsPerPage = res.hitsPerPage()
	res.NbPages = res.nbPages()
	res.Page = res.page()
	return res
}

func (res *Response) Header() http.Header {
	h := make(http.Header)
	h.Set(headers.ContentType, param.NdJSON)

	return h
}

func (res *Response) FilterByFacetValue(attr, val string) []map[string]any {
	f, err := res.facets.GetFacet(attr)
	if err != nil {
		return res.results
	}
	items := lo.ToAnySlice(f.FindByValue(val).RelatedTo)
	return FilterDataByID(res.results, items, res.UID)
}

func (res *Response) nbHits() int {
	return len(res.results)
}

func (res *Response) nbPages() int {
	hpp := res.HitsPerPage

	if hpp < 1 {
		return 1
	}

	nb := res.NbHits / hpp
	if r := res.NbHits % hpp; r > 0 {
		nb++
	}

	return nb
}

func (res *Response) hitsPerPage() int {
	if res.Params.Has(param.HitsPerPage) {
		return res.Params.HitsPerPage
	}
	return viper.GetInt(param.HitsPerPage.Snake())
}

func (res *Response) page() int {
	if !res.Params.Has(param.Page) {
		return 0
	}
	return res.Params.Page
}

func (res *Response) visibleHits() []map[string]any {
	if res.HitsPerPage == -1 ||
		res.NbHits < res.HitsPerPage {
		return res.results
	}

	if nb := res.NbPages; res.Page >= nb {
		return []map[string]any{}
	}

	return lo.Subset(res.results, res.Page*res.HitsPerPage, uint(res.HitsPerPage))
}

func (res *Response) FilterResults() []map[string]any {
	if res.Facets == nil {
		return res.results
	}

	var hits []map[string]any
	for idx, d := range res.results {
		if i, ok := d[res.UID]; ok {
			idx = cast.ToInt(i)
		}
		if res.facets.Bitmap().ContainsInt(idx) {
			hits = append(hits, d)
		}
	}
	return hits
}
