package srch

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

const (
	Hits                  = `hits`
	SearchableAttributes  = `searchableAttributes`
	AttributesForFaceting = `attributesForFaceting`
	AttributesToRetrieve  = `attributesToRetrieve`
	Page                  = "page"
	HitsPerPage           = "hitsPerPage"
	SortFacetValuesBy     = `sortFacetValuesBy`
	ParamQuery            = `query`
	ParamFacets           = "facets"
	ParamFacetFilters     = `facetFilters`
	ParamFilters          = "filters"
	DataDir               = `dataDir`
	DataFile              = `dataFile`
	ParamFullText         = `fullText`
	NbHits                = `nbHits`
	NbPages               = `nbPages`
	DefaultField          = `title`
)

type Query struct {
	Params url.Values
}

func NewQuery(q any) *Query {
	return &Query{
		Params: ParseQuery(q),
	}
}

func ParseQuery(queries ...any) url.Values {
	q := make(url.Values)
	for _, query := range queries {
		vals, err := ParseValues(query)
		if err != nil {
			continue
		}
		for k, val := range vals {
			for _, v := range val {
				q.Add(k, v)
			}
		}
	}
	return q
}

func (q *Query) Merge(queries ...*Query) {
	for _, query := range queries {
		for k, val := range query.Params {
			for _, v := range val {
				q.Params.Add(k, v)
			}
		}
	}
}

func (q Query) GetData() ([]map[string]any, error) {
	if !q.HasData() {
		return nil, errors.New("no data")
	}
	return GetDataFromQuery(&q.Params)
}

func (q Query) HasData() bool {
	return q.Params.Has(DataFile) || q.Params.Has(DataDir)
}

func (q Query) GetSettings() *Settings {
	s := defaultSettings()
	s.setValsFromQuery(&q)
	return s
}

func (q Query) GetFacetFilters() (*Filters, error) {
	if !q.HasFilters() {
		return nil, errors.New("no filters")
	}
	f, err := DecodeFilter(q.Params.Get(ParamFacetFilters))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (q Query) Query() string {
	return q.Params.Get(ParamQuery)
}

func (q *Query) Fields() []*Field {
	attrs := q.GetSrchAttr()
	fields := make([]*Field, len(attrs))
	for i, attr := range attrs {
		fields[i] = NewField(attr, q.GetAnalyzer())
	}
	return fields
}

func (q Query) Facets() []*Field {
	facets := q.FacetAttrs()
	fields := make([]*Field, len(facets))
	for i, attr := range facets {
		fields[i] = NewField(attr, OrFacet)
	}
	return fields
}

func (q Query) FacetAttrs() []string {
	if !q.Params.Has(ParamFacets) {
		q.Params[ParamFacets] = q.GetFacetAttr()
	}
	return q.Params[ParamFacets]
}

func (q Query) Page() int {
	p := q.Params.Get(Page)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}

func (q Query) HitsPerPage() int {
	p := q.Params.Get(HitsPerPage)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}

func (q Query) HasFilters() bool {
	return q.Params.Has(ParamFacetFilters)
}

func (q Query) Get(key string) []string {
	if q.Params.Has(key) {
		return q.Params[key]
	}
	return []string{}
}

func (q Query) GetSrchAttr() []string {
	return GetQueryStringSlice(SrchAttr, q.Params)
}

func (q Query) GetFacetAttr() []string {
	return GetQueryStringSlice(FacetAttr, q.Params)
}

func (q Query) GetAnalyzer() string {
	return GetAnalyzer(q.Params)
}

func (q Query) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(q.Encode())
	if err != nil {
		return nil, err
	}
	return d, err
}

func (q *Query) UnmarshalJSON(d []byte) error {
	var p string
	err := json.Unmarshal(d, &p)
	if err != nil {
		return err
	}

	err = q.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

func (q Query) Encode() string {
	return q.Params.Encode()
}

func (q *Query) Decode(str string) error {
	var err error
	q.Params, err = url.ParseQuery(str)
	return err
}

func GetDataFromQuery(q *url.Values) ([]map[string]any, error) {
	var data []map[string]any
	var err error
	switch {
	case q.Has(DataFile):
		qu := *q
		data, err = FileSrc(qu[DataFile]...)
		q.Del(DataFile)
	case q.Has(DataDir):
		data, err = DirSrc(q.Get(DataDir))
		q.Del(DataDir)
	}
	return data, err
}

func GetQueryStringSlice(key string, q url.Values) []string {
	var vals []string
	if q.Has(key) {
		for _, val := range q[key] {
			if val == "" {
				break
			}
			for _, v := range strings.Split(val, ",") {
				vals = append(vals, v)
			}
		}
	}
	if key == SrchAttr {
		switch len(vals) {
		case 0:
			return []string{DefaultField}
		case 1:
			if vals[0] == "" {
				return []string{DefaultField}
			}
		}
	}
	return vals
}

// ParseValues takes an interface{} and returns a url.Values.
func ParseValues(f any) (url.Values, error) {
	filters := make(url.Values)
	var err error
	switch val := f.(type) {
	case url.Values:
		return val, nil
	case []byte:
		return ParseQueryBytes(val)
	case string:
		return ParseQueryString(val)
	default:
		filters, err = cast.ToStringMapStringSliceE(val)
		if err != nil {
			return nil, err
		}
	}
	return filters, nil
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, err
}

// ParseQueryBytes parses a byte slice to url.Values.
func ParseQueryBytes(val []byte) (url.Values, error) {
	filters, err := cast.ToStringMapStringSliceE(string(val))
	if err != nil {
		return nil, err
	}
	return url.Values(filters), err
}
