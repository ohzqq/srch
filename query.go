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
	Hits                 = `hits`
	AttributesToRetrieve = `attributesToRetrieve`
	SrchAttr             = `searchableAttributes`
	FacetAttr            = `attributesForFaceting`
	Page                 = "page"
	HitsPerPage          = "hitsPerPage"
	SortFacetsBy         = `sortFacetValuesBy`
	Query                = `query`
	ParamFacets          = "facets"
	ParamFilters         = "filters"
	FacetFilters         = `facetFilters`
	DataDir              = `dataDir`
	DataFile             = `dataFile`
	ParamFullText        = `fullText`
	NbHits               = `nbHits`
	NbPages              = `nbPages`
	DefaultField         = `title`
)

type Params struct {
	Values url.Values
}

func NewQuery(q any) *Params {
	return &Params{
		Values: ParseQuery(q),
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

func GetAnalyzer(q url.Values) string {
	if q.Has(ParamFullText) {
		return Text
	}
	return Fuzzy
}

func (q *Params) Merge(queries ...*Params) {
	for _, query := range queries {
		for k, val := range query.Values {
			for _, v := range val {
				q.Values.Add(k, v)
			}
		}
	}
}

func (q Params) GetData() ([]map[string]any, error) {
	if !q.HasData() {
		return nil, errors.New("no data")
	}
	return GetDataFromQuery(&q.Values)
}

func (q Params) HasData() bool {
	return q.Values.Has(DataFile) || q.Values.Has(DataDir)
}

func (q Params) GetFacetFilters() (*Filters, error) {
	if !q.HasFilters() {
		return nil, errors.New("no filters")
	}
	f, err := DecodeFilter(q.Values.Get(FacetFilters))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (q Params) Query() string {
	return q.Values.Get(Query)
}

func (q Params) GetSrchAttr() []string {
	return GetQueryStringSlice(SrchAttr, q.Values)
}

func (p *Params) SortFacetsBy() string {
	sort := "count"
	if p.Values.Has(SortFacetsBy) {
		if by := p.Values.Get(SortFacetsBy); by == "count" || by == "alpha" {
			sort = by
		}
	}
	return sort
}

func (q *Params) SrchAttr() []string {
	if !q.Values.Has(SrchAttr) {
		return []string{DefaultField}
	}
	q.Values[SrchAttr] = q.GetSrchAttr()
	if len(q.Values[SrchAttr]) < 1 {
		q.Values[SrchAttr] = []string{DefaultField}
	}
	return q.Values[SrchAttr]
}

func (q *Params) Fields() []*Field {
	attrs := q.SrchAttr()
	fields := make([]*Field, len(attrs))
	for i, attr := range attrs {
		fields[i] = NewTextField(attr, q)
	}
	return fields
}

func (q *Params) Facets() []*Field {
	facets := q.FacetAttr()
	fields := make([]*Field, len(facets))
	for i, attr := range facets {
		fields[i] = NewFacet(attr, q)
	}
	return fields
}

func (q Params) FacetAttr() []string {
	if !q.Values.Has(ParamFacets) {
		q.Values[ParamFacets] = q.GetFacetAttr()
	}
	return q.Values[ParamFacets]
}

func (q Params) GetFacetAttr() []string {
	return GetQueryStringSlice(FacetAttr, q.Values)
}

func (q Params) Page() int {
	p := q.Values.Get(Page)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}

func (q Params) HitsPerPage() int {
	p := q.Values.Get(HitsPerPage)
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return page
}

func (q Params) HasFilters() bool {
	return q.Values.Has(FacetFilters)
}

func (q Params) Get(key string) []string {
	if q.Values.Has(key) {
		return q.Values[key]
	}
	return []string{}
}

func (q Params) GetAnalyzer() string {
	return GetAnalyzer(q.Values)
}

func (q Params) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(q.Encode())
	if err != nil {
		return nil, err
	}
	return d, err
}

func (q *Params) UnmarshalJSON(d []byte) error {
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

func (q Params) Encode() string {
	return q.Values.Encode()
}

func (q Params) String() string {
	return q.Values.Encode()
}

func (q *Params) Decode(str string) error {
	var err error
	q.Values, err = url.ParseQuery(str)
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
