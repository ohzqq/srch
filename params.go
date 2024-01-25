package srch

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/ohzqq/srch/txt"
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
	SortBy               = `sortBy`
	Order                = `order`
)

type Params struct {
	Values url.Values
}

func NewQuery(q any) *Params {
	return &Params{
		Values: ParseQuery(q),
	}
}

func (p Params) GetSlice(key string) []string {
	if p.Values.Has(key) {
		return p.Values[key]
	}
	return []string{}
}

func (p Params) Get(key string) string {
	if p.Values.Has(key) {
		return p.Values.Get(key)
	}
	return ""
}

func (p *Params) Merge(queries ...*Params) {
	for _, query := range queries {
		for k, val := range query.Values {
			for _, v := range val {
				p.Values.Add(k, v)
			}
		}
	}
}

func (p *Params) NewField(attr string) *Field {
	f := NewField(attr)
	if p.IsFullText() {
		f.Tokens = txt.NewTokens(txt.Fulltext())
	}
	f.SortBy = p.SortFacetsBy()
	return f
}

func (p *Params) newFields(attrs []string) []*Field {
	fields := make([]*Field, len(attrs))
	for i, attr := range attrs {
		fields[i] = p.NewField(attr)
	}
	return fields
}

func (p *Params) SrchAttr() []string {
	if !p.Values.Has(SrchAttr) {
		return []string{DefaultField}
	}
	p.Values[SrchAttr] = GetQueryStringSlice(SrchAttr, p.Values)
	if len(p.Values[SrchAttr]) < 1 {
		p.Values[SrchAttr] = []string{DefaultField}
	}
	return p.Values[SrchAttr]
}

func (p *Params) Fields() []*Field {
	return p.newFields(p.SrchAttr())
}

func (p Params) HasFilters() bool {
	return p.Values.Has(FacetFilters)
}

func (p Params) GetFacetFilters() (*Filters, error) {
	if !p.HasFilters() {
		return nil, errors.New("no filters")
	}
	f, err := DecodeFilter(p.Values.Get(FacetFilters))
	if err != nil {
		return nil, err
	}
	return f, nil
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

func (p Params) FacetAttr() []string {
	if !p.Values.Has(ParamFacets) {
		p.Values[ParamFacets] = GetQueryStringSlice(FacetAttr, p.Values)
	}
	return p.Values[ParamFacets]
}

func (p *Params) Facets() []*Field {
	return p.newFields(p.FacetAttr())
}

func (p Params) Page() int {
	pn := p.Values.Get(Page)
	page, err := strconv.Atoi(pn)
	if err != nil {
		return 0
	}
	return page
}

func (p Params) HitsPerPage() int {
	pn := p.Values.Get(HitsPerPage)
	page, err := strconv.Atoi(pn)
	if err != nil {
		return 0
	}
	return page
}

func (p *Params) IsFullText() bool {
	return p.Values.Has(ParamFullText)
}

func (p Params) Query() string {
	return p.Values.Get(Query)
}

func (p Params) GetAnalyzer() string {
	if p.Values.Has(ParamFullText) {
		return Text
	}
	return Keyword
}

func (p Params) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(p.Encode())
	if err != nil {
		return nil, err
	}
	return d, err
}

func (p *Params) UnmarshalJSON(d []byte) error {
	var pn string
	err := json.Unmarshal(d, &pn)
	if err != nil {
		return err
	}

	err = p.Decode(pn)
	if err != nil {
		return err
	}
	return nil
}

func (p Params) Encode() string {
	return p.Values.Encode()
}

func (p Params) String() string {
	return p.Values.Encode()
}

func (p *Params) Decode(str string) error {
	var err error
	p.Values, err = url.ParseQuery(str)
	return err
}

func (p Params) HasData() bool {
	return p.Values.Has(DataFile) || p.Values.Has(DataDir)
}

func (p Params) GetData() ([]map[string]any, error) {
	if !p.HasData() {
		return nil, errors.New("no data")
	}
	return GetDataFromQuery(&p.Values)
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
